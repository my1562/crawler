package apiclient

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/my1562/crawler/models"
)

type ApiClient struct {
	client *resty.Client
}

func New(client *resty.Client) *ApiClient {
	return &ApiClient{client: client}
}

type AddressResponse struct {
	Result *models.AddressAr `json:"result,omitempty"`
}

func (api *ApiClient) TakeNextAddress() (interface{}, error) {
	resp, err := api.client.R().
		SetResult(&AddressResponse{}).
		Post("/address/take")
	if err != nil {
		return nil, err
	}
	fmt.Println(string(resp.Body()))
	result := resp.Result()
	return result, nil
}
