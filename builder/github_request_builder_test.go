package builder

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/LasramR/sclng-backend-test-lasramR/model/version"
)

func TestGithubRequestBuilder_UnsupportedVersion(t *testing.T) {
	unsupportedVersion := version.GithubAPIVersion("Unsupported")
	_, err := NewGithubRequestBuilder(unsupportedVersion)

	if err == nil {
		t.Fatalf("Github API version %s should NOT be supported", unsupportedVersion)
	}
}

func TestGithubRequestBuilder_20221128(t *testing.T) {
	grb, err := NewGithubRequestBuilder(version.GITHUB_API_2022_11_28)

	if err != nil {
		t.Fatalf("Github API version %s should be supported", version.GITHUB_API_2022_11_28)
	}

	_ = grb.With("language", "Python")
	_ = grb.With("license", "apache-2.0")
	_ = grb.Limit(80)

	req, err := grb.Build(context.Background(), http.MethodGet, "/search/repositories")

	if err != nil {
		t.Fatalf("Github request builder with valid parameters should not return an error")
	}

	if !strings.Contains(req.URL.String(), url.QueryEscape("is:public")) {
		t.Fatalf("GithubRequestBuilder %s missing is:public from built request URL", version.GITHUB_API_2022_11_28)
	}

	if !strings.Contains(req.URL.String(), url.QueryEscape("language:Python")) {
		t.Fatalf("GithubRequestBuilder %s missing language:Python from built request URL", version.GITHUB_API_2022_11_28)
	}

	if !strings.Contains(req.URL.String(), url.QueryEscape("license:apache-2.0")) {
		t.Fatalf("GithubRequestBuilder %s missing license:apache-2.0 from built request URL", version.GITHUB_API_2022_11_28)
	}

	if !strings.Contains(req.URL.String(), "per_page=80") {
		t.Fatalf("GithubRequestBuilder %s missing per_page=80 from built request URL", version.GITHUB_API_2022_11_28)
	}

	if !strings.Contains(req.URL.String(), "page=1") {
		t.Fatalf("GithubRequestBuilder %s missing page=1 from built request URL", version.GITHUB_API_2022_11_28)
	}
}
