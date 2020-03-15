package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/my1562/crawler/apiclient"
	"github.com/my1562/crawler/config"
	"github.com/my1562/geocoder"
	"go.uber.org/dig"
)

func main() {

	c := dig.New()
	c.Provide(func(conf *config.Config) *resty.Client {
		client := resty.New().SetHostURL(conf.ApiUrl)
		return client
	})
	c.Provide(config.NewConfig)
	c.Provide(apiclient.New)
	c.Provide(
		func(conf *config.Config) (*geocoder.Geocoder, error) {
			//TODO: add config param
			geo := geocoder.NewGeocoder("./data/gobs/geocoder-data.gob")
			geo.BuildSpatialIndex(100)
			return geo, nil
		})

	err := c.Invoke(func(client *apiclient.ApiClient, geo *geocoder.Geocoder) {
		addr, err := client.TakeNextAddress()
		if err != nil {
			panic(err)
		}
		addressAr := geo.AddressByID(uint32(addr.ID))
		if addressAr == nil {
			panic("No such address") //TODO: skip it and notify invalid address
		}
		fmt.Printf("(%d) %s %s", addressAr.Street1562.ID, addressAr.Street1562.Name, addressAr.Address.GetBuildingAsString())
	})
	if err != nil {
		panic(err)
	}

}
