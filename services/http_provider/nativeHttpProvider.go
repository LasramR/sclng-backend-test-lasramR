package http_provider

import (
	"encoding/json"
	"io"
	"net/http"
)

type NativeHttpClient struct {
	Do         func(req *http.Request) (*http.Response, error)
	NewRequest func(method string, url string, body io.Reader) (*http.Request, error)
}

type NativeHttpProvider struct {
	client *NativeHttpClient
}

func (provider *NativeHttpProvider) GetJson(url string, payload any) error {
	request, err := provider.client.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	response, err := provider.client.Do(request)

	if err != nil {
		return err
	}

	return json.NewDecoder(response.Body).Decode(payload)
}

func NewNativeHttpProvider(client *NativeHttpClient) *NativeHttpProvider {
	return &NativeHttpProvider{
		client: client,
	}
}
