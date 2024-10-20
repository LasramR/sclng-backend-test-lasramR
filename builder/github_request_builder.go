package builder

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

type GithubRequestBuilder interface {
	Build(ctx context.Context, method, baseUrl string) (*http.Request, error)
	Authorization(value string)
	With(key, value string)
	Sort(value string)
}

type GithubRequestBuilder20221128 struct {
	authorization string
	q             map[string]string
	sort          string
}

func (grb *GithubRequestBuilder20221128) Build(ctx context.Context, method, baseUrl string) (*http.Request, error) {
	hrb := NewHttpRequestBuilder(method, baseUrl)

	if grb.authorization != "" {
		hrb.AddHeader("Authorization", []string{fmt.Sprintf("Bearer %s", grb.authorization)})
	}

	if len(grb.q) != 0 {
		stringifiedQ := make([]string, 0, len(grb.q))
		for k, v := range grb.q {
			stringifiedQ = append(stringifiedQ, fmt.Sprintf("%s:%s", k, v))
		}
		hrb.AddQueryParam("q", strings.Join(stringifiedQ, " "))
	}

	if grb.sort != "" {
		hrb.AddQueryParam("sort", grb.sort)
	}

	return hrb.BuildRequest(ctx)
}

func (grb *GithubRequestBuilder20221128) Authorization(value string) {
	if value != "" {
		grb.authorization = value
	}
}

func (grb *GithubRequestBuilder20221128) With(key, value string) {
	grb.q[key] = value
}

var sortValues = []string{"created", "updated", "comments"}

func (grb *GithubRequestBuilder20221128) Sort(value string) {
	if slices.Contains(sortValues, value) {
		grb.sort = value
	}
}

func NewGithubRequestBuilder20221128() *GithubRequestBuilder20221128 {
	return &GithubRequestBuilder20221128{
		authorization: "",
		q:             make(map[string]string),
		sort:          "",
	}
}
