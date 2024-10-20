package builder

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestGithubRequestBuilder20221128(t *testing.T) {
	grb := NewGithubRequestBuilder20221128()

	grb.With("language", "Python")
	grb.With("license", "apache2.0")
	grb.Sort("something")

	req, err := grb.Build(context.Background(), http.MethodGet, "https://api.github.com")

	expectedUrl := fmt.Sprintf("https://api.github.com?q=%s", url.QueryEscape("language:Python license:apache2.0"))

	if err != nil {
		t.Fatalf("Github request builder with valid parameters should not return an error")
	}

	if req.URL.String() != expectedUrl {

		t.Fatalf("Github request builder 2022-11-28 formating error")
	}
}
