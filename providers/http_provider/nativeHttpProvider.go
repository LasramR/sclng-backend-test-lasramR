package http_provider

import (
	"context"
	"encoding/json"
	"errors"
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

func (provider *NativeHttpProvider) GetJson(ctx context.Context, url string, payload any) error {
	request, err := provider.client.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	responseCh := make(chan error)
	defer close(responseCh)

	go func() {
		response, err := provider.client.Do(request)
		if err != nil {
			responseCh <- err
		} else {
			defer response.Body.Close()
			responseCh <- json.NewDecoder(response.Body).Decode(payload)
		}
	}()
	select {
	case err := <-responseCh:
		return err
	case <-ctx.Done():
		return errors.New("request timeout")
	}
}

func NewNativeHttpProvider(client *NativeHttpClient) *NativeHttpProvider {
	return &NativeHttpProvider{
		client: client,
	}
}
