package main

import (
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/my1562/crawler/apiclient"
	"github.com/my1562/crawler/config"
	"github.com/my1562/crawler/tasks"
	"github.com/my1562/crawler/worker"
	"go.uber.org/dig"
)

func main() {

	c := dig.New()
	c.Provide(func(conf *config.Config) *resty.Client {
		client := resty.New().SetHostURL(conf.ApiUrl)
		return client
	})
	c.Provide(config.NewConfig)
	c.Provide(apiclient.NewApiClient)
	c.Provide(tasks.NewTasks)
	c.Provide(worker.NewWorker)

	err := c.Invoke(func(tasks *tasks.Tasks, worker *worker.Worker) {
		go func() {
			if err := worker.Listen(); err != nil {
				log.Panic(err)
			}
		}()

		retryInterval := time.Second * 10

		for {
			delay, err := tasks.GetDelay()
			if err != nil {
				log.Printf("GetDelay error %s\n", err)
				log.Printf("Retrying in %d seconds\n", retryInterval/time.Second)
				time.Sleep(retryInterval)
				continue
			}
			err = tasks.GetNextAddressCheckAndStore()
			if err != nil {
				log.Printf("GetNextAddressCheckAndStore error %s\n", err)
				log.Printf("Retrying in %d seconds\n", retryInterval/time.Second)
				time.Sleep(retryInterval)
				continue
			}

			log.Printf("Sleeping %f minutes\n", float32(delay)/float32(time.Minute))
			time.Sleep(delay)
		}
	})
	if err != nil {
		panic(err)
	}
}
