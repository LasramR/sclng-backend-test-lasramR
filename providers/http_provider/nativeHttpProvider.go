package http_provider

import (
	"encoding/json"
	"net/http"
)

type NativeHttpClient struct {
	Do func(req *http.Request) (*http.Response, error)
}

type NativeHttpProvider struct {
	client NativeHttpClient
}

func (provider *NativeHttpProvider) ReqUnmarshalledBody(req *http.Request, unMarshalledResBody any) error {
	response, err := provider.client.Do(req)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(unMarshalledResBody)
}

func NewNativeHttpProvider(client NativeHttpClient) *NativeHttpProvider {
	return &NativeHttpProvider{
		client: client,
	}
}
