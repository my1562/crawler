package tasks

import (
	"errors"
	"fmt"
	"time"

	my1562client "github.com/my1562/client"
	"github.com/my1562/crawler/apiclient"
	"github.com/my1562/crawler/models"
	"github.com/my1562/crawler/utils"
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

func (tasks *Tasks) GetDelay() (time.Duration, error) {
	addrCount, err := tasks.client.GetAddressCount()
	if err != nil {
		return 0, err
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

	return time.Duration(delay) * time.Second, nil
}

func (tasks *Tasks) GetNextAddressCheckAndStore() error {
	addr, err := tasks.client.TakeNextAddress()
	if err != nil {
		return err
	}
	fullAddress := tasks.geo.AddressByID(uint32(addr.ID))
	if fullAddress == nil {
		return errors.New("No such address") //TODO: skip it and notify invalid address
	}
	building := fullAddress.Address.GetBuildingAsString()
	streetID := int(fullAddress.Street1562.ID)
	fmt.Printf("(Address.ID=%d) (Street1562.ID=%d) %s %s\n", fullAddress.Address.ID, streetID, fullAddress.Street1562.Name, building)
	status, err := my1562client.GetStatus(streetID, building)
	if err != nil {
		return err
	}

	message := utils.FormatServiceMessage(status)
	fmt.Println("Status:")
	fmt.Println(message)
	fmt.Printf("HasMessage:  %t\n", status.HasMessage)

	var checkStatus models.AddressArCheckStatus
	if status.HasMessage {
		checkStatus = models.AddressStatusWork
	} else {
		checkStatus = models.AddressStatusNoWork
	}

	err = tasks.client.UpdateAddress(int64(fullAddress.Address.ID), checkStatus, message, status.Hash)
	if err != nil {
		return err
	}

	return nil
}
