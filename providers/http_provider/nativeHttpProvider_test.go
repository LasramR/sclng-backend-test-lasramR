package http_provider

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type DoReturnT struct {
	response *http.Response
	err      error
}

type NewRequestReturnT struct {
	request *http.Request
	err     error
}

func MockHttpClient(DoReturn DoReturnT, NewRequestReturn NewRequestReturnT) *NativeHttpClient {
	return &NativeHttpClient{
		Do: func(req *http.Request) (*http.Response, error) {
			return DoReturn.response, DoReturn.err
		},
		NewRequest: func(method string, url string, body io.Reader) (*http.Request, error) {
			return NewRequestReturn.request, NewRequestReturn.err
		},
	}
}

type GetJsonResponseT struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

func TestGetJson_NewRequestFails(t *testing.T) {
	httpProvider := NewNativeHttpProvider(MockHttpClient(
		DoReturnT{nil, nil},
		NewRequestReturnT{nil, errors.New("NewRequest")},
	))

	var result GetJsonResponseT
	err := httpProvider.GetJson(context.Background(), "https://somedataendpoint.io", &result)

	if err.Error() != "NewRequest" {
		t.Fatalf("should return an error when NewRequest fails")
	}
}

func TestGetJson_DoFails(t *testing.T) {
	httpProvider := NewNativeHttpProvider(MockHttpClient(
		DoReturnT{nil, errors.New("Do")},
		NewRequestReturnT{&http.Request{}, nil},
	))

	var result GetJsonResponseT
	err := httpProvider.GetJson(context.Background(), "https://somedataendpoint.io", &result)

	if err.Error() != "Do" {
		t.Fatalf("should return an error when Do fails")
	}
}

func TestGetJson_DecodeFails(t *testing.T) {
	httpProvider := NewNativeHttpProvider(MockHttpClient(
		DoReturnT{&http.Response{
			Body: io.NopCloser(bytes.NewReader([]byte{})),
		}, nil},
		NewRequestReturnT{&http.Request{}, nil},
	))

	var result GetJsonResponseT
	err := httpProvider.GetJson(context.Background(), "https://somedataendpoint.io", &result)

	if err == nil {
		t.Fatalf("should return an error when unmarshalling fails")
	}
}

func TestGetJson_Timeout(t *testing.T) {
	httpProvider := NewNativeHttpProvider(&NativeHttpClient{
		Do: func(req *http.Request) (*http.Response, error) {
			time.Sleep(time.Second * 10)
			return nil, nil
		},
		NewRequest: func(method string, url string, body io.Reader) (*http.Request, error) {
			return &http.Request{}, nil
		},
	})

	timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), time.Second*1)
	defer cancelTimeout()

	var result GetJsonResponseT
	err := httpProvider.GetJson(timeoutCtx, "https://somedataendpoint.io", &result)

	if err.Error() != "request timeout" {
		t.Fatalf("should have failed because of context timeout")
	}
}

func TestGetJson_Valid(t *testing.T) {
	httpProvider := NewNativeHttpProvider(MockHttpClient(
		DoReturnT{&http.Response{
			Body: io.NopCloser(bytes.NewReader([]byte("{\"key\":\"a\", \"value\":1}"))),
		}, nil},
		NewRequestReturnT{&http.Request{}, nil},
	))

	var expected = GetJsonResponseT{Key: "a", Value: 1}

	var result GetJsonResponseT
	err := httpProvider.GetJson(context.Background(), "https://somedataendpoint.io", &result)

	if err != nil {
		t.Fatalf("should not have return an error")
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("should have fed our result struct")
	}
}
