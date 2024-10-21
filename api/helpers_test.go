package api

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/LasramR/sclng-backend-test-lasramR/builder"
	"github.com/LasramR/sclng-backend-test-lasramR/model"
	"github.com/LasramR/sclng-backend-test-lasramR/providers"
	"github.com/LasramR/sclng-backend-test-lasramR/repositories"
	"github.com/LasramR/sclng-backend-test-lasramR/util"
	"github.com/redis/go-redis/v9"
)

type MockGitHubService struct {
	err error
}

func (mgs MockGitHubService) GetGithubProjectsWithStats(ctx context.Context, grb builder.GithubRequestBuilder) (repositories.GithubRepositoriesResult, error) {
	if mgs.err != nil {
		return repositories.GithubRepositoriesResult{}, mgs.err
	}

	return repositories.GithubRepositoriesResult{
		Total:            1,
		IncompleteResult: false,
		Repositories: []*model.Repository{
			{
				FullName:    "fmuiin14/BlazingTool",
				Owner:       "fmuiin14",
				Description: "Brute force ethereum wallet mnemonics",
				Repository:  "fmuiin14/BlazingTool",
				License: util.NullableJsonField[string]{
					Value:  "mit",
					IsNull: false,
				},
				RepositoryUrl: "https://github.com/fmuiin14/BlazingTool",
				CreatedAt:     "2024-10-19T10:17:16Z",
				UpdatedAt:     "2024-10-20T16:36:13Z",
				Languages: model.Language{
					"JavaScript": model.LanguageStats{
						Bytes: 1548,
					},
					"SCSS": model.LanguageStats{
						Bytes: 250,
					},
				},
				Size: 156464,
			},
		},
	}, nil
}

func MochCacheProvider(value string, getErr, setErr error) providers.CacheProvider {
	return providers.NewRedisCacheProvider(&providers.RedisClient{
		Get: func(ctx context.Context, s string) *redis.StringCmd {
			cmd := &redis.StringCmd{}

			if getErr != nil {
				cmd.SetErr(getErr)
			} else {
				cmd.SetVal(value)
			}

			return cmd
		},
		Set: func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
			cmd := &redis.StatusCmd{}

			cmd.SetErr(setErr)

			return cmd
		},
	})
}

type MockResponseWriter struct {
	StatusCode int
	Buffer     *bytes.Buffer
}

func (mrw *MockResponseWriter) Write(b []byte) (int, error) {
	return mrw.Buffer.Write(b)
}

func (mrw *MockResponseWriter) WriteHeader(statusCode int) {
	mrw.StatusCode = statusCode
}

func (mrw *MockResponseWriter) Header() http.Header {
	return http.Header{}
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{
		Buffer: &bytes.Buffer{},
	}
}
