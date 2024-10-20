package builder

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestGithubRequestBuilder_UnsupportedVersion(t *testing.T) {
	unsupportedVersion := GithubAPIVersion("Unsupported")
	_, err := NewGithubRequestBuilder(unsupportedVersion)

	if err == nil {
		t.Fatalf("Github API version %s should NOT be supported", unsupportedVersion)
	}
}

func TestGithubRequestBuilder_20221128(t *testing.T) {
	grb, err := NewGithubRequestBuilder(GITHUB_API_2022_11_28)

	if err != nil {
		t.Fatalf("Github API version %s should be supported", GITHUB_API_2022_11_28)
	}

	grb.With("language", "Python")
	grb.With("license", "apache2.0")

	req, err := grb.Build(context.Background(), http.MethodGet, "https://api.github.com")

	expectedUrl := fmt.Sprintf("https://api.github.com?q=%s", url.QueryEscape("language:Python license:apache2.0"))

	if err != nil {
		t.Fatalf("Github request builder with valid parameters should not return an error")
	}

	if req.URL.String() != expectedUrl {

		t.Fatalf("GithubRequestBuilder %s URL formating error %s", GITHUB_API_2022_11_28, req.URL.String())
	}
}
