package tasks

import (
	"fmt"
	"time"

	my1562client "github.com/my1562/client"
	"github.com/my1562/crawler/apiclient"
	"github.com/my1562/crawler/models"
	"github.com/my1562/geocoder"
)

const (
	BestTotalTime int64 = 2 * 60 * 60
	MinDelay            = 5 * 60
	MaxDelay            = 30 * 60
)

type Tasks struct {
	client *apiclient.ApiClient
	geo    *geocoder.Geocoder
}

func New(
	client *apiclient.ApiClient,
	geo *geocoder.Geocoder,
) *Tasks {
	return &Tasks{
		client: client,
		geo:    geo,
	}
}

func (tasks *Tasks) GetDelay() time.Duration {
	addrCount, err := tasks.client.GetAddressCount()
	if err != nil {
		panic(err)
	}
	if addrCount <= 0 {
		addrCount = 1
	}
	delay := BestTotalTime / addrCount
	if delay > MaxDelay {
		delay = MaxDelay
	}
	if delay < MinDelay {
		delay = MinDelay
	}

	return time.Duration(delay) * time.Second
}

func (tasks *Tasks) GetNextAddressCheckAndStore() {
	addr, err := tasks.client.TakeNextAddress()
	if err != nil {
		panic(err)
	}
	fullAddress := tasks.geo.AddressByID(uint32(addr.ID))
	if fullAddress == nil {
		panic("No such address") //TODO: skip it and notify invalid address
	}
	building := fullAddress.Address.GetBuildingAsString()
	streetID := int(fullAddress.Street1562.ID)
	fmt.Printf("(Address.ID=%d) (Street1562.ID=%d) %s %s\n", fullAddress.Address.ID, streetID, fullAddress.Street1562.Name, building)
	status, err := my1562client.GetStatus(streetID, building)
	if err != nil {
		panic(err)
	}
	fmt.Println("Status:")
	fmt.Printf(" - Title       %s\n", status.Title)
	fmt.Printf(" - Description %s\n", status.Description)
	fmt.Printf(" - HasMessage  %t\n", status.HasMessage)

	var message string
	var checkStatus models.AddressArCheckStatus
	if status.HasMessage {
		checkStatus = models.AddressStatusWork
		message = status.Title + "\n" + status.Description
	} else {
		checkStatus = models.AddressStatusNoWork
		message = ""
	}

	err = tasks.client.UpdateAddress(int64(fullAddress.Address.ID), checkStatus, message)
	if err != nil {
		panic(err)
	}

}
