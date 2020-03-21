package apiclient

import (
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/my1562/crawler/models"
)

type ApiClient struct {
	client *resty.Client
}

type ErrorResponse struct {
	ErrorValue string `json:"error,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.ErrorValue
}

func New(client *resty.Client) *ApiClient {
	client.SetDebug(true)
	return &ApiClient{client: client}
}

func (api *ApiClient) TakeNextAddress() (*models.AddressAr, error) {

	type AddressResponse struct {
		Result *models.AddressAr `json:"result,omitempty"`
	}

	resp, err := api.client.R().
		SetResult(&AddressResponse{}).
		SetError(&ErrorResponse{}).
		Post("/address/take")
	if err != nil {
		return nil, err
	}
	if backendError := resp.Error(); backendError != nil {
		return nil, backendError.(*ErrorResponse)
	}
	result := resp.Result()
	return result.(*AddressResponse).Result, nil
}

func (api *ApiClient) GetAddressCount() (int64, error) {
	type Response struct {
		Result int64
	}

	resp, err := api.client.R().
		SetResult(&Response{}).
		SetError(&ErrorResponse{}).
		Get("/address/count")
	if err != nil {
		return 0, err
	}
	if backendError := resp.Error(); backendError != nil {
		return 0, backendError.(*ErrorResponse)
	}
	result := resp.Result()
	return result.(*Response).Result, nil
}

func (api *ApiClient) UpdateAddress(id int64, status models.AddressArCheckStatus, message string, hash string) error {
	type UpdateAddressRequest struct {
		CheckStatus    models.AddressArCheckStatus
		ServiceMessage string
		Hash           string
	}

	resp, err := api.client.R().
		SetPathParams(map[string]string{
			"id": strconv.FormatInt(id, 10),
		}).
		SetError(&ErrorResponse{}).
		SetBody(&UpdateAddressRequest{
			CheckStatus:    status,
			ServiceMessage: message,
			Hash:           hash,
		}).
		Put("/address/{id}")
	if err != nil {
		return err
	}
	if backendError := resp.Error(); backendError != nil {
		return backendError.(*ErrorResponse)
	}

	return nil
}
