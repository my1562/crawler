package apiclient

import (
	"github.com/go-resty/resty/v2"
	"github.com/my1562/crawler/models"
)

type ApiClient struct {
	client *resty.Client
}

func New(client *resty.Client) *ApiClient {
	return &ApiClient{client: client}
}

func (api *ApiClient) TakeNextAddress() (*models.AddressAr, error) {

	type AddressResponse struct {
		Result *models.AddressAr `json:"result,omitempty"`
	}

	resp, err := api.client.R().
		SetResult(&AddressResponse{}).
		Post("/address/take")
	if err != nil {
		return nil, err
	}
	result := resp.Result()
	return result.(*AddressResponse).Result, nil
}
