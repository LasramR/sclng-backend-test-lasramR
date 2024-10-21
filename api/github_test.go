package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/LasramR/sclng-backend-test-lasramR/model"
	"github.com/LasramR/sclng-backend-test-lasramR/model/version"
	"github.com/LasramR/sclng-backend-test-lasramR/util"
)

func TestGitHubProjectsHandler_WrongMethod(t *testing.T) {
	handler := GitHubProjectsHandler(
		MockGitHubService{},
		MochCacheProvider("", errors.New("not in cache"), nil),
		5,
		version.GITHUB_API_2022_11_28,
	)

	r, _ := http.NewRequest(http.MethodPost, "http://endpoint.io", nil)
	w := NewMockResponseWriter()

	err := handler(w, r, nil)

	if err != nil {
		t.Fatalf("api handler should not return an error")
	}

	if w.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("only http GET is allowed on this endpoint")
	}
}

func TestGitHubProjectsHandler_Valid(t *testing.T) {
	mgs := MockGitHubService{}
	handler := GitHubProjectsHandler(
		&mgs,
		MochCacheProvider("", errors.New("not in cache"), nil),
		5,
		version.GITHUB_API_2022_11_28,
	)

	r, _ := http.NewRequest(http.MethodGet, "http://endpoint.io", nil)
	w := NewMockResponseWriter()

	err := handler(w, r, nil)

	if err != nil {
		t.Fatalf("api handler should not return an error")
	}

	result, _ := mgs.GetGithubProjectsWithStats(r.Context(), nil)
	expected, _ := json.Marshal(model.ApiListResponse[[]*model.Repository]{
		TotalCount:       1,
		Count:            1,
		Content:          result.Repositories,
		IncompleteResult: false,
		Next:             util.NextFullUrlFromRequest(r),
		Previous: util.NullableJsonField[string]{
			IsNull: true,
			Value:  "",
		},
	})

	if w.StatusCode != http.StatusOK {
		t.Fatalf("Should have responded with status 200")
	}

	if !reflect.DeepEqual(w.Buffer.Bytes()[:len(w.Buffer.Bytes())-1], expected) {
		t.Fatalf("Expected should have been written in response writter")
	}
}

func TestGitHubProjectsHandler_UnvalidLimit(t *testing.T) {
	mgs := MockGitHubService{}
	handler := GitHubProjectsHandler(
		&mgs,
		MochCacheProvider("", errors.New("not in cache"), nil),
		5,
		version.GITHUB_API_2022_11_28,
	)

	r, _ := http.NewRequest(http.MethodGet, "http://endpoint.io?limit=pouet", nil)
	w := NewMockResponseWriter()

	err := handler(w, r, nil)

	if err != nil {
		t.Fatalf("api handler should not return an error")
	}

	if w.StatusCode != http.StatusBadRequest {
		t.Fatalf("Should have responded with status 400")
	}
}

func TestGitHubProjectsHandler_UnvalidPage(t *testing.T) {
	mgs := MockGitHubService{}
	handler := GitHubProjectsHandler(
		&mgs,
		MochCacheProvider("", errors.New("not in cache"), nil),
		5,
		version.GITHUB_API_2022_11_28,
	)

	r, _ := http.NewRequest(http.MethodGet, "http://endpoint.io?page=teuop", nil)
	w := NewMockResponseWriter()

	err := handler(w, r, nil)

	if err != nil {
		t.Fatalf("api handler should not return an error")
	}

	if w.StatusCode != http.StatusBadRequest {
		t.Fatalf("Should have responded with status 400")
	}
}

func TestGitHubProjectsHandler_UnsupportedQueryParam(t *testing.T) {
	mgs := MockGitHubService{}
	handler := GitHubProjectsHandler(
		&mgs,
		MochCacheProvider("", errors.New("not in cache"), nil),
		5,
		version.GITHUB_API_2022_11_28,
	)

	r, _ := http.NewRequest(http.MethodGet, "http://endpoint.io?unsupported=ohno", nil)
	w := NewMockResponseWriter()

	err := handler(w, r, nil)

	if err != nil {
		t.Fatalf("api handler should not return an error")
	}

	if w.StatusCode != http.StatusBadRequest {
		t.Fatalf("Should have responded with status 400")
	}
}

func TestGitHubProjectsHandler_UnsupportedGithubApiVersion(t *testing.T) {
	mgs := MockGitHubService{}
	handler := GitHubProjectsHandler(
		&mgs,
		MochCacheProvider("", errors.New("not in cache"), nil),
		5,
		version.GithubAPIVersion("unsupported"),
	)

	r, _ := http.NewRequest(http.MethodGet, "http://endpoint.io?page=teuop", nil)
	w := NewMockResponseWriter()

	err := handler(w, r, nil)

	if err != nil {
		t.Fatalf("api handler should not return an error")
	}

	if w.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("Should have responded with status StatusServiceUnavailable")
	}
}

func TestGitHubProjectsHandler_GithubServiceError(t *testing.T) {
	mgs := MockGitHubService{err: errors.New("request error :/")}
	handler := GitHubProjectsHandler(
		&mgs,
		MochCacheProvider("", errors.New("not in cache"), nil),
		5,
		version.GITHUB_API_2022_11_28,
	)

	r, _ := http.NewRequest(http.MethodGet, "http://endpoint.io", nil)
	w := NewMockResponseWriter()

	err := handler(w, r, nil)

	if err != nil {
		t.Fatalf("api handler should not return an error")
	}

	if w.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Should have responded with status StatusInternalServerError")
	}
}
