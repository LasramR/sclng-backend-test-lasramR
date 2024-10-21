package builder

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/LasramR/sclng-backend-test-lasramR/model/version"
	"github.com/LasramR/sclng-backend-test-lasramR/util"
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

// Following types are used for function composition in order to abstract the request building process

type githubParamSetter func(hrb *util.HttpRequestBuilder, params map[string]string)
type githubSortSetter func(hrb *util.HttpRequestBuilder, sort string)
type authorizationSetter func(hrb *util.HttpRequestBuilder, authorization string)
type limitSetter func(hrb *util.HttpRequestBuilder, limit int)
type pageSetter func(hrb *util.HttpRequestBuilder, page int)

type githubRequestBuilderAPIVersionned struct {
	apiVersion              version.GithubAPIVersion
	apiBaseUrl              string
	authorizationSetterFunc authorizationSetter
	authorizationValue      string
	supportedParams         []string
	paramSetterFunc         githubParamSetter
	params                  map[string]string
	supportedSort           []string
	sortSetter              githubSortSetter
	sortBy                  string
	limitSetterFunc         limitSetter
	maxLimit                int
	limitValue              int
	pageSetterFunc          pageSetter
	pageValue               int
}

func (grb *githubRequestBuilderAPIVersionned) Build(ctx context.Context, method, url string) (*http.Request, error) {
	var fullUrl string
	if strings.HasSuffix(grb.apiBaseUrl, "/") {
		fullUrl = grb.apiBaseUrl + url
	} else if grb.apiBaseUrl == "" {
		fullUrl = url
	} else {
		fullUrl = fmt.Sprintf("%s%s", grb.apiBaseUrl, url)
	}

	hrb := util.NewHttpRequestBuilder(method, fullUrl)

	if grb.authorizationValue != "" {
		grb.authorizationSetterFunc(hrb, grb.authorizationValue)
	}

	if len(grb.params) != 0 {
		grb.paramSetterFunc(hrb, grb.params)
	}

	if grb.sortBy != "" {
		grb.sortSetter(hrb, grb.sortBy)
	}

	grb.limitSetterFunc(hrb, grb.limitValue)
	grb.pageSetterFunc(hrb, grb.pageValue)

	return hrb.BuildRequest(ctx)
}

func (grb *githubRequestBuilderAPIVersionned) Authorization(value string) {
	if value != "" {
		grb.authorizationValue = value
	}
}
func (grb *githubRequestBuilderAPIVersionned) With(key, value string) error {
	if slices.Contains(grb.supportedParams, key) && value != "" {
		grb.params[key] = value
		return nil
	}

	return fmt.Errorf("%s parameter is not supported", key)
}
func (grb *githubRequestBuilderAPIVersionned) Sort(value string) error {
	if slices.Contains(grb.supportedSort, value) {
		grb.sortBy = value
		return nil
	}

	return fmt.Errorf("%s sorting is not supported [%s] allowed", value, strings.Join(grb.supportedSort, ","))
}

func (grb *githubRequestBuilderAPIVersionned) Limit(value int) error {
	if value < 1 || grb.maxLimit < value {
		return fmt.Errorf("parameter limit %d is exceeding max limit of %d", value, grb.maxLimit)
	}

	grb.limitValue = value
	return nil
}

func (grb *githubRequestBuilderAPIVersionned) Page(value int) error {
	if value < 1 {
		return fmt.Errorf("page parameter %d must be greater than 0", value)
	}

	grb.pageValue = value
	return nil
}

// Factory method that creates a GithubRequestBuilder for a specific API version, err != nil if API version is not supported
func NewGithubRequestBuilder(ApiVersion version.GithubAPIVersion) (GithubRequestBuilder, error) {
	switch ApiVersion {
	case version.GITHUB_API_2022_11_28:
		return &githubRequestBuilderAPIVersionned{
			apiVersion:    version.GITHUB_API_2022_11_28,
			apiBaseUrl:    "https://api.github.com",
			supportedSort: []string{"updated", "forks", "stars"},
			authorizationSetterFunc: func(hrb *util.HttpRequestBuilder, authorization string) {
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
			paramSetterFunc: func(hrb *util.HttpRequestBuilder, params map[string]string) {
				stringifiedParams := make([]string, 0, len(params))
				for _, k := range util.SortedKeys(params) { // Ensure that request url is deterministic for caching purposes
					v := params[k]
					stringifiedParams = append(stringifiedParams, fmt.Sprintf("%s:%s", k, v))
				}
				hrb.AddQueryParam("q", strings.Join(stringifiedParams, " "))
			},
			sortSetter: func(hrb *util.HttpRequestBuilder, sort string) {
				hrb.AddQueryParam("sort", sort)
			},
			maxLimit: 100,
			limitSetterFunc: func(hrb *util.HttpRequestBuilder, limit int) {
				hrb.AddQueryParam("per_page", fmt.Sprintf("%d", limit))
			},
			limitValue: 100,
			pageSetterFunc: func(hrb *util.HttpRequestBuilder, page int) {
				hrb.AddQueryParam("page", fmt.Sprintf("%d", page))
			},
			pageValue: 1,
		}, nil
	default:
		return nil, errors.New("unsupported github api version")
	}
}
