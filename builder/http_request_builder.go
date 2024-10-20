package builder

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type HttpRequestBuilder struct {
	BaseUrl     string
	Method      string
	queryParams url.Values
	headers     http.Header
}

func (hrb *HttpRequestBuilder) AddQueryParam(key, value string) {
	hrb.queryParams.Add(key, value)
}

func (hrb *HttpRequestBuilder) AddHeader(key string, value []string) {
	hrb.headers[key] = value
}

func (hrb *HttpRequestBuilder) BuildRequest(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, hrb.Method, fmt.Sprintf("%s?%s", hrb.BaseUrl, hrb.queryParams.Encode()), nil)

	if err != nil {
		return nil, err
	}

	req.Header = hrb.headers

	return req, nil
}

func NewHttpRequestBuilder(method, baseUrl string) *HttpRequestBuilder {
	return &HttpRequestBuilder{
		BaseUrl:     baseUrl,
		Method:      method,
		queryParams: url.Values{},
		headers:     make(map[string][]string),
	}
}
