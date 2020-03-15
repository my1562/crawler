package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/my1562/crawler/apiclient"
	"github.com/my1562/crawler/config"
	"github.com/my1562/crawler/tasks"
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
	c.Provide(tasks.New)
	c.Provide(
		func(conf *config.Config) (*geocoder.Geocoder, error) {
			//TODO: add config param
			geo := geocoder.NewGeocoder("./data/gobs/geocoder-data.gob")
			geo.BuildSpatialIndex(100)
			return geo, nil
		})

	err := c.Invoke(func(tasks *tasks.Tasks) {
		tasks.GetNextAddressCheckAndStore()
	})
	if err != nil {
		panic(err)
	}

}
