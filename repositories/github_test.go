package repositories

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/LasramR/sclng-backend-test-lasramR/builder"
	"github.com/LasramR/sclng-backend-test-lasramR/model"
	"github.com/LasramR/sclng-backend-test-lasramR/model/version"
	"github.com/LasramR/sclng-backend-test-lasramR/providers"
	"github.com/LasramR/sclng-backend-test-lasramR/util"
	"github.com/redis/go-redis/v9"
)

func MockHttpProvider(bodys []string, errs []error) providers.HttpProvider {
	i := 0
	return providers.NewNativeHttpProvider(
		providers.NativeHttpClient{
			Do: func(req *http.Request) (*http.Response, error) {
				errIdx := len(errs) - 1
				if i <= len(errs)-1 {
					errIdx = i
				}
				if errs != nil && len(errs) <= errIdx && errs[errIdx] != nil {
					return nil, errs[errIdx]
				}

				bodyIdx := len(bodys) - 1

				if i <= len(bodys)-1 {
					bodyIdx = i
				}

				resp := &http.Response{
					Body: io.NopCloser(bytes.NewReader([]byte(bodys[bodyIdx]))),
				}

				i += 1

				return resp, nil
			},
		},
	)
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

func TestGetGithubProjectsWithStats_API20221128(t *testing.T) {
	gr, err := NewGithubApiRepository(
		version.GITHUB_API_2022_11_28,
		MockHttpProvider(
			[]string{GITHUB_SEARCH_REPOS_RESPONSE_BODY_SAMPLE, GITHUB_LANGUAGE_RESPONSE_BODY_SAMPLE_1, GITHUB_LANGUAGE_RESPONSE_BODY_SAMPLE_2},
			nil,
		),
		MochCacheProvider("", errors.New("no value in cache"), nil),
		5,
		"sometoken",
	)

	if err != nil {
		t.Fatalf("github repository should support API %s", version.GITHUB_API_2022_11_28)
	}

	grb, _ := builder.NewGithubRequestBuilder(version.GITHUB_API_2022_11_28)

	result, _ := gr.GetManyRepositories(context.Background(), grb)

	expected := &GithubRepositoriesResult{
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
	}

	if len(result.Repositories) != 2 {
		t.Fatalf("Should have returned two record")
	}

	if !reflect.DeepEqual(expected.Repositories[0], expected.Repositories[0]) {
		t.Fatalf("Should equals expected value %v", expected.Repositories[0])
	}
}

const (
	GITHUB_SEARCH_REPOS_RESPONSE_BODY_SAMPLE = `
{
  "total_count": 606814,
  "incomplete_results": false,
  "items": [
    {
      "id": 875186168,
      "name": "BlazingTool",
      "full_name": "fmuiin14/BlazingTool",
      "owner": {
        "login": "fmuiin14"
      },
      "size": 156464,
      "description": "Brute force ethereum wallet mnemonics",
      "languages_url": "https://api.github.com/repos/fmuiin14/BlazingTool/languages",
      "created_at": "2024-10-19T10:17:16Z",
      "updated_at": "2024-10-20T16:36:13Z",
      "license": {
        "key": "mit",
        "name": "MIT License"
      }
    },
    {
      "id": 875187065,
      "name": "ShadowTool",
      "full_name": "fmuiin14/ShadowTool",
      "owner": {
        "login": "fmuiin14"
      },
      "size": 9856,
      "description": "This script is designed to automatically generate seed phrases and check balances for Tron networks. If a wallet with a non-zero balance is found, the wallet's information (address, mnemonic, private key, and balances) is logged and saved to a file named result.txt.",
      "languages_url": "https://api.github.com/repos/fmuiin14/ShadowTool/languages",
      "created_at": "2024-10-19T10:20:11Z",
      "updated_at": "2024-10-20T16:36:05Z",
      "license": {
        "key": "mit",
        "name": "MIT License"
      }
    }
  ]
}`
	GITHUB_LANGUAGE_RESPONSE_BODY_SAMPLE_1 = `{
			"JavaScript": 1548,
			"SCSS": 250
		}`

	GITHUB_LANGUAGE_RESPONSE_BODY_SAMPLE_2 = `{
			"TypeScript": 10884,
			"Bash": 170
		}`
)
