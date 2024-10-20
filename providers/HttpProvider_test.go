package provider

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/LasramR/sclng-backend-test-lasramR/util"
)

func MockHttpClient(DoReturn util.Result[*http.Response]) NativeHttpClient {
	return NativeHttpClient{
		Do: func(req *http.Request) (*http.Response, error) {
			return DoReturn.Value, DoReturn.Error
		},
	}
}

type GetUnmarshalledResponseT struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

func TestReqUnmarshalledBody_DoFails(t *testing.T) {
	httpProvider := NewNativeHttpProvider(MockHttpClient(
		util.Result[*http.Response]{Value: nil, Error: errors.New("Do")},
	))

	req, _ := http.NewRequest(http.MethodGet, "https://somedataendpoint.io", nil)
	var result GetUnmarshalledResponseT

	err := httpProvider.ReqUnmarshalledBody(req, &result)

	if err.Error() != "Do" {
		t.Fatalf("should return an error when Do fails")
	}
}

func TestReqUnmarshalledBody_DecodeFails(t *testing.T) {
	httpProvider := NewNativeHttpProvider(MockHttpClient(
		util.Result[*http.Response]{
			Value: &http.Response{
				Body: io.NopCloser(bytes.NewReader([]byte{})),
			},
			Error: nil,
		},
	))

	req, _ := http.NewRequest(http.MethodGet, "https://somedataendpoint.io", nil)
	var result GetUnmarshalledResponseT

	err := httpProvider.ReqUnmarshalledBody(req, &result)

	if err == nil {
		t.Fatalf("should return an error when unmarshalling fails")
	}
}

func TestGetJson_Valid(t *testing.T) {
	httpProvider := NewNativeHttpProvider(MockHttpClient(
		util.Result[*http.Response]{
			Value: &http.Response{
				Body: io.NopCloser(bytes.NewReader([]byte("{\"key\":\"a\", \"value\":1}"))),
			},
			Error: nil,
		},
	))

	var expected = GetUnmarshalledResponseT{Key: "a", Value: 1}

	req, _ := http.NewRequest(http.MethodGet, "https://somedataendpoint.io", nil)
	var result GetUnmarshalledResponseT

	err := httpProvider.ReqUnmarshalledBody(req, &result)

	if err != nil {
		t.Fatalf("should not have return an error")
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("should have fed our result struct")
	}
}
