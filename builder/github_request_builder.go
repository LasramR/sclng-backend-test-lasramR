package builder

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

// Provide mecanisms to build a Github API request
type GithubRequestBuilder interface {
	Build(ctx context.Context, method, baseUrl string) (*http.Request, error)
	Authorization(value string)
	With(key, value string) error
	Sort(value string) error
}

type GithubAPIVersion string

// Supported github api version
const (
	GITHUB_API_2022_11_28 GithubAPIVersion = "2022-11-28"
)

type GithubParamSetter func(hrb *HttpRequestBuilder, params map[string]string)
type GithubSortSetter func(hrb *HttpRequestBuilder, sort string)
type AuthorizationSetter func(hrb *HttpRequestBuilder, authorization string)

type GithubRequestBuilderAPIVersionned struct {
	ApiVersion              GithubAPIVersion
	ApiBaseUrl              string
	authorizationSetterFunc AuthorizationSetter
	authorization           string
	supportedParams         []string
	paramSetterFunc         GithubParamSetter
	params                  map[string]string
	supportedSort           []string
	sortSetter              GithubSortSetter
	sortBy                  string
}

func (grb *GithubRequestBuilderAPIVersionned) Build(ctx context.Context, method, url string) (*http.Request, error) {
	var fullUrl string
	if strings.HasSuffix(grb.ApiBaseUrl, "/") {
		fullUrl = grb.ApiBaseUrl + url
	} else if grb.ApiBaseUrl == "" {
		fullUrl = url
	} else {
		fullUrl = fmt.Sprintf("%s%s", grb.ApiBaseUrl, url)
	}

	hrb := NewHttpRequestBuilder(method, fullUrl)

	if grb.authorization != "" {
		grb.authorizationSetterFunc(hrb, grb.authorization)
	}

	if len(grb.params) != 0 {
		grb.paramSetterFunc(hrb, grb.params)
	}

	if grb.sortBy != "" {
		grb.sortSetter(hrb, grb.sortBy)
	}

	return hrb.BuildRequest(ctx)
}

func (grb *GithubRequestBuilderAPIVersionned) Authorization(value string) {
	if value != "" {
		grb.authorization = value
	}
}
func (grb *GithubRequestBuilderAPIVersionned) With(key, value string) error {
	if slices.Contains(grb.supportedParams, key) && value != "" {
		grb.params[key] = value
		return nil
	}

	return fmt.Errorf("%s parameter is not supported", key)
}
func (grb *GithubRequestBuilderAPIVersionned) Sort(value string) error {
	if slices.Contains(grb.supportedSort, value) {
		grb.sortBy = value
		return nil
	}

	return fmt.Errorf("%s sorting is not supported", value)
}

func NewGithubRequestBuilder(ApiVersion GithubAPIVersion) (GithubRequestBuilder, error) {
	switch ApiVersion {
	case GITHUB_API_2022_11_28:
		return &GithubRequestBuilderAPIVersionned{
			ApiVersion:    GITHUB_API_2022_11_28,
			ApiBaseUrl:    "https://api.github.com",
			supportedSort: []string{"created", "updated", "comments"},
			authorizationSetterFunc: func(hrb *HttpRequestBuilder, authorization string) {
				hrb.AddHeader("Authorization", []string{fmt.Sprintf("Bearer %s", authorization)})
			},
			supportedParams: []string{
				"language",
				"license",
				"user",
				"org",
				"repo",
			},
			params: map[string]string{"is": "public"},
			paramSetterFunc: func(hrb *HttpRequestBuilder, params map[string]string) {
				stringifiedParams := make([]string, 0, len(params))
				for k, v := range params {
					stringifiedParams = append(stringifiedParams, fmt.Sprintf("%s:%s", k, v))
				}
				hrb.AddQueryParam("q", strings.Join(stringifiedParams, " "))
			},
			sortSetter: func(hrb *HttpRequestBuilder, sort string) {
				// TODO handle sorting
			},
		}, nil
	default:
		return nil, errors.New("unsupported github api version")
	}
}
