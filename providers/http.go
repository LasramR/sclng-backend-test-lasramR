package providers

import (
	"encoding/json"
	"net/http"
)

// Allow to perform http related operations
type HttpProvider interface {
	// Perform a HTTP request and unmarshals the response body into unMarshalledResBody argument
	ReqUnmarshalledBody(req *http.Request, unMarshalledResBody any) error
}

// IoC of the http client
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

func NewNativeHttpProvider(client NativeHttpClient) HttpProvider {
	return &NativeHttpProvider{
		client: client,
	}
}
