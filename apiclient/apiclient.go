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

func NewApiClient(client *resty.Client) *ApiClient {
	client.SetDebug(true)
	return &ApiClient{client: client}
}

type ShortGeocoderAddress struct {
	ID             uint32
	Address        string
	Building       string
	Street1562ID   uint32
	Street1562Name string
}
type TakeNextResponse struct {
	Address         *models.AddressAr
	GeocoderAddress *ShortGeocoderAddress
}

func (api *ApiClient) TakeNextAddress(addressID int64) (*TakeNextResponse, error) {

	type AddressResponse struct {
		Result *TakeNextResponse `json:"result,omitempty"`
	}

	var url string
	request := api.client.R().
		SetResult(&AddressResponse{}).
		SetError(&ErrorResponse{})
	if addressID == 0 {
		url = "/address-take"
	} else {
		url = "/address-take/{id}"
		request = request.
			SetPathParams(map[string]string{
				"id": strconv.FormatInt(addressID, 10),
			})
	}

	resp, err := request.Post(url)
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
		Get("/address-count")
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
