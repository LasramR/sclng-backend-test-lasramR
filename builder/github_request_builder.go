package builder

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

// Provide mecanisms to build complex request to query Github API request
type GithubRequestBuilder interface {
	// Build the request object
	Build(ctx context.Context, method, baseUrl string) (*http.Request, error)
	// Attach an authorization header
	Authorization(value string)
	// Adds a query parameter, error != nil if parameter "key" is not supported
	With(key, value string) error
	// Adds a sort parameter, error != nil if sorting "value" is not supported
	Sort(value string) error
	// Limits a request result count by "value", error != nil if "value" is invalid
	Limit(value int) error
	// Limits a request result count by "value", error != nil if "value" is invalid
	Page(value int) error
}

type GithubAPIVersion string

// Supported github api version by the GithubRequestBuilder
const (
	GITHUB_API_2022_11_28 GithubAPIVersion = "2022-11-28"
)

type githubParamSetter func(hrb *HttpRequestBuilder, params map[string]string)
type githubSortSetter func(hrb *HttpRequestBuilder, sort string)
type authorizationSetter func(hrb *HttpRequestBuilder, authorization string)
type limitSetter func(hrb *HttpRequestBuilder, limit int)
type pageSetter func(hrb *HttpRequestBuilder, page int)

type githubRequestBuilderAPIVersionned struct {
	ApiVersion              GithubAPIVersion
	ApiBaseUrl              string
	AuthorizationSetterFunc authorizationSetter
	AuthorizationValue      string
	SupportedParams         []string
	ParamSetterFunc         githubParamSetter
	Params                  map[string]string
	SupportedSort           []string
	SortSetter              githubSortSetter
	SortBy                  string
	LimitSetterFunc         limitSetter
	MaxLimit                int
	LimitValue              int
	PageSetterFunc          pageSetter
	PageValue               int
}

func (grb *githubRequestBuilderAPIVersionned) Build(ctx context.Context, method, url string) (*http.Request, error) {
	var fullUrl string
	if strings.HasSuffix(grb.ApiBaseUrl, "/") {
		fullUrl = grb.ApiBaseUrl + url
	} else if grb.ApiBaseUrl == "" {
		fullUrl = url
	} else {
		fullUrl = fmt.Sprintf("%s%s", grb.ApiBaseUrl, url)
	}

	hrb := NewHttpRequestBuilder(method, fullUrl)

	if grb.AuthorizationValue != "" {
		grb.AuthorizationSetterFunc(hrb, grb.AuthorizationValue)
	}

	if len(grb.Params) != 0 {
		grb.ParamSetterFunc(hrb, grb.Params)
	}

	if grb.SortBy != "" {
		grb.SortSetter(hrb, grb.SortBy)
	}

	grb.LimitSetterFunc(hrb, grb.LimitValue)
	grb.PageSetterFunc(hrb, grb.PageValue)

	return hrb.BuildRequest(ctx)
}

func (grb *githubRequestBuilderAPIVersionned) Authorization(value string) {
	if value != "" {
		grb.AuthorizationValue = value
	}
}
func (grb *githubRequestBuilderAPIVersionned) With(key, value string) error {
	if slices.Contains(grb.SupportedParams, key) && value != "" {
		grb.Params[key] = value
		return nil
	}

	return fmt.Errorf("%s parameter is not supported", key)
}
func (grb *githubRequestBuilderAPIVersionned) Sort(value string) error {
	if slices.Contains(grb.SupportedSort, value) {
		grb.SortBy = value
		return nil
	}

	return fmt.Errorf("%s sorting is not supported", value)
}

func (grb *githubRequestBuilderAPIVersionned) Limit(value int) error {
	if value < 1 || grb.MaxLimit < value {
		return fmt.Errorf("parameter limit %d is exceeding max limit of %d", value, grb.MaxLimit)
	}

	grb.LimitValue = value
	return nil
}

func (grb *githubRequestBuilderAPIVersionned) Page(value int) error {
	if value < 1 {
		return fmt.Errorf("page parameter %d must be greater than 0", value)
	}

	grb.PageValue = value
	return nil
}

// Factory method that creates a GithubRequestBuilder for a specific API version, err != nil if API version is not supported
func NewGithubRequestBuilder(ApiVersion GithubAPIVersion) (GithubRequestBuilder, error) {
	switch ApiVersion {
	case GITHUB_API_2022_11_28:
		return &githubRequestBuilderAPIVersionned{
			ApiVersion:    GITHUB_API_2022_11_28,
			ApiBaseUrl:    "https://api.github.com",
			SupportedSort: []string{"created", "updated", "comments"},
			AuthorizationSetterFunc: func(hrb *HttpRequestBuilder, authorization string) {
				hrb.AddHeader("Authorization", []string{fmt.Sprintf("Bearer %s", authorization)})
			},
			SupportedParams: []string{
				"language",
				"license",
				"user",
				"org",
				"repo",
			},
			Params: map[string]string{"is": "public"},
			ParamSetterFunc: func(hrb *HttpRequestBuilder, params map[string]string) {
				stringifiedParams := make([]string, 0, len(params))
				for k, v := range params {
					stringifiedParams = append(stringifiedParams, fmt.Sprintf("%s:%s", k, v))
				}
				hrb.AddQueryParam("q", strings.Join(stringifiedParams, " "))
			},
			SortSetter: func(hrb *HttpRequestBuilder, sort string) {
				// TODO handle sorting
			},
			MaxLimit: 100,
			LimitSetterFunc: func(hrb *HttpRequestBuilder, limit int) {
				hrb.AddQueryParam("per_page", fmt.Sprintf("%d", limit))
			},
			LimitValue: 100,
			PageSetterFunc: func(hrb *HttpRequestBuilder, page int) {
				hrb.AddQueryParam("page", fmt.Sprintf("%d", page))
			},
			PageValue: 1,
		}, nil
	default:
		return nil, errors.New("unsupported github api version")
	}
}
