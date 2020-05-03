package tasks

import (
	"log"
	"time"

	my1562client "github.com/my1562/client"
	"github.com/my1562/crawler/apiclient"
	"github.com/my1562/crawler/models"
	"github.com/my1562/crawler/utils"
)

const (
	BestTotalTime int64 = 2 * 60 * 60
	MinDelay            = 5 * 60
	MaxDelay            = 30 * 60
)

type Tasks struct {
	client *apiclient.ApiClient
}

func NewTasks(
	client *apiclient.ApiClient,
) *Tasks {
	return &Tasks{
		client: client,
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
	log.Printf("Regular address check")
	return tasks.getAddressCheckAndStore(0)
}

func (tasks *Tasks) PriorityCheck(addressID int64) error {
	log.Printf("Prioritized address check: %d", addressID)
	return tasks.getAddressCheckAndStore(addressID)
}

func (tasks *Tasks) getAddressCheckAndStore(addressID int64) error {
	log.Println("---- getAddressCheckAndStore START ----")
	defer log.Println("---- getAddressCheckAndStore END ----")

	addr, err := tasks.client.TakeNextAddress(addressID)
	if err != nil {
		return err
	}

	id := addr.Address.ID
	building := addr.GeocoderAddress.Building
	streetID := int(addr.GeocoderAddress.Street1562ID)

	log.Printf("(streetID=%d) %s\n", streetID, addr.GeocoderAddress.Address)

	status, err := my1562client.GetStatus(streetID, building)
	if err != nil {
		return err
	}

	message := utils.FormatServiceMessage(status)
	log.Println("Status:")
	log.Println(message)
	log.Printf("HasMessage:  %t\n", status.HasMessage)

	var checkStatus models.AddressArCheckStatus
	if status.HasMessage {
		checkStatus = models.AddressStatusWork
	} else {
		checkStatus = models.AddressStatusNoWork
	}

	err = tasks.client.UpdateAddress(id, checkStatus, message, status.Hash)
	if err != nil {
		return err
	}

	return nil
}
