package util

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewHttpRequestBuilder_New(t *testing.T) {
	expected := &HttpRequestBuilder{
		Method:      http.MethodGet,
		BaseUrl:     "https://somesite.io",
		queryParams: url.Values{},
		headers:     make(map[string][]string),
	}

	actual := NewHttpRequestBuilder(http.MethodGet, "https://somesite.io")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Invalid HttpRequestBuilder")
	}
}

func TestAddQueryParam(t *testing.T) {
	requestBuilder := NewHttpRequestBuilder(http.MethodGet, "https://somesite.io")

	requestBuilder.AddQueryParam("a", "123")

	expected := url.Values{"a": {"123"}}

	if !reflect.DeepEqual(expected, requestBuilder.queryParams) {
		t.Fatalf("AddQueryParam should modify internal queryParam property")
	}
}

func TestAddHeader(t *testing.T) {
	requestBuilder := NewHttpRequestBuilder(http.MethodGet, "https://somesite.io")

	requestBuilder.AddHeader("Authorization", []string{"SECRET"})

	expected := map[string][]string{
		"Authorization": {"SECRET"},
	}

	if !reflect.DeepEqual(expected["Authorization"], requestBuilder.headers["Authorization"]) {
		t.Fatalf("AddHeader should modify internal headers property")
	}
}

func TestBuildRequest_InvalidRequest(t *testing.T) {
	requestBuilder := NewHttpRequestBuilder(http.MethodGet, "https://somesite.io")

	_, err := requestBuilder.BuildRequest(context.Context(nil))

	if err == nil {
		t.Fatalf("invalid BuildRequest should return an error")
	}
}

func TestBuildRequest_ParametizedURL(t *testing.T) {
	requestBuilder := NewHttpRequestBuilder(http.MethodGet, "https://somesite.io")
	requestBuilder.AddQueryParam("a", "123")

	expected := "https://somesite.io?a=123"
	request, err := requestBuilder.BuildRequest(context.Background())

	if err != nil {
		t.Fatalf("valid BuildRequest should not return an error")
	}

	actual := request.URL.String()

	if actual != expected {
		t.Fatalf("QueryParam of a RequestBuilder should be encoded in request url ")
	}
}

func TestBuildRequest_CustomHeaders(t *testing.T) {
	requestBuilder := NewHttpRequestBuilder(http.MethodGet, "https://somesite.io")
	requestBuilder.AddHeader("Authorization", []string{"SECRET"})

	expected := map[string][]string{
		"Authorization": {"SECRET"},
	}
	request, err := requestBuilder.BuildRequest(context.Background())

	if err != nil {
		t.Fatalf("valid BuildRequest should not return an error")
	}

	actual := request.Header

	if reflect.DeepEqual(actual, expected) {
		t.Fatalf("")
	}
}
