package util

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Return the complete URL from a query including protocol and host with the given url parameters
// Edited from https://gist.github.com/karl-gustav/001e05e70527986f8b6d11f675ed610c
func fullUrlFrom(r *http.Request, queryParams url.Values) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s%s?%s#%s", scheme, r.Host, r.URL.Path, queryParams.Encode(), r.URL.Fragment)
}

// Return the complete URL from a query including protocol and host
func FullUrlFromRequest(r *http.Request) string {
	return fullUrlFrom(r, r.URL.Query())
}

// Given a request, returns its url with query param page incremented or set to 2
func NextFullUrlFromRequest(r *http.Request) string {
	queryParams := r.URL.Query()

	if page, err := strconv.Atoi(queryParams.Get("page")); err == nil {
		queryParams.Set("page", fmt.Sprintf("%d", page+1))
	} else {
		queryParams.Set("page", "2")
	}

	return fullUrlFrom(r, queryParams)
}

// Given a request, returns its url with query param page decremented.
// Will return an empty string if page=1
func PreviousFullUrlFromRequest(r *http.Request) string {
	queryParams := r.URL.Query()

	if page, err := strconv.Atoi(queryParams.Get("page")); err != nil {
		return ""
	} else {
		queryParams.Set("page", fmt.Sprintf("%d", page-1))
	}

	return fullUrlFrom(r, queryParams)
}
