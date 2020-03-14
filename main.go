package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/my1562/crawler/apiclient"
	"github.com/my1562/crawler/config"
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

	err := c.Invoke(func(client *apiclient.ApiClient) {
		client.TakeNextAddress()
	})
	if err != nil {
		panic(err)
	}

}
