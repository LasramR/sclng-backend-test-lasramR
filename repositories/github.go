package repositories

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/LasramR/sclng-backend-test-lasramR/builder"
	"github.com/LasramR/sclng-backend-test-lasramR/model"
	"github.com/LasramR/sclng-backend-test-lasramR/model/external"
	"github.com/LasramR/sclng-backend-test-lasramR/model/version"
	"github.com/LasramR/sclng-backend-test-lasramR/providers"
	"github.com/LasramR/sclng-backend-test-lasramR/util"
)

// Represent an aggregated Response from Github API
type GithubRepositoriesResult struct {
	Repositories []*model.Repository `json:"repositories"`
	// Describe the total number of matching projects from the response api, not len(Repositories)
	Total int `json:"total"`
	// Set to true if some sub aggregations failed
	IncompleteResult bool `json:"incomplete_result"`
}

// Allow to interact with the GitHub REST API
type GithubApiRepository interface {
	// Fetch many repositories, error != nil if
	GetManyRepositories(ctx context.Context, grb builder.GithubRequestBuilder) (GithubRepositoriesResult, error)
}

// Parametized implementation of the GitHub repository that abstracts the entity mapping process
type githubVersionnedApiRepository[T any, M util.Mappable[T]] struct {
	githubToken        string
	httpProvider       providers.HttpProvider
	cacheProvider      providers.CacheProvider
	mapperFunc         util.MapperFunc[T, *model.Repository]
	cacheDurationInMin time.Duration
}

func (gr *githubVersionnedApiRepository[T, M]) GetManyRepositories(ctx context.Context, grb builder.GithubRequestBuilder) (GithubRepositoriesResult, error) {
	grb.Authorization(gr.githubToken)
	req, err := grb.Build(ctx, http.MethodGet, "/search/repositories")

	if err != nil {
		return GithubRepositoriesResult{}, nil
	}

	requestUrl := req.URL.String()
	var repositories GithubRepositoriesResult

	if err = gr.cacheProvider.GetUnmarshalled(ctx, requestUrl, &repositories); err == nil {
		return repositories, nil
	}

	var apiResponse M
	err = gr.httpProvider.ReqUnmarshalledBody(req, &apiResponse)

	if err != nil {
		return GithubRepositoriesResult{}, err
	}

	mapped, errorsCollected := util.AsyncListMapper(
		ctx,
		apiResponse,
		gr.mapperFunc,
		time.Second*30,
	)

	// If we had some results, we cache it
	if len(errorsCollected) != len(apiResponse.Items()) {
		_ = gr.cacheProvider.SetMarshalled(ctx, requestUrl, repositories, time.Minute*gr.cacheDurationInMin)
	}

	return GithubRepositoriesResult{
		Repositories:     mapped,
		Total:            apiResponse.Count(),
		IncompleteResult: len(errorsCollected) != 0,
	}, nil
}

// Factory method that creates a GithubApiRepository for a specific API version, err != nil if API version is not supported
func NewGithubApiRepository(apiVersion version.GithubAPIVersion, httpProvider providers.HttpProvider, cacheProvider providers.CacheProvider, cacheDurationInMin time.Duration, githubToken string) (GithubApiRepository, error) {
	switch apiVersion {
	case version.GITHUB_API_2022_11_28:
		return &githubVersionnedApiRepository[external.RepositoriesResponseItem, external.RepositoriesResponse]{
			githubToken:        githubToken,
			httpProvider:       httpProvider,
			cacheProvider:      cacheProvider,
			cacheDurationInMin: cacheDurationInMin,
			// Mapper function converts items of the external model to our model
			mapperFunc: func(ctx context.Context, rawRepository external.RepositoriesResponseItem) (*model.Repository, error) {
				req, err := http.NewRequest(http.MethodGet, rawRepository.LanguagesUrl, nil)

				if err != nil {
					return nil, nil
				}

				requestUrl := rawRepository.LanguagesUrl
				var repository model.Repository

				if err = cacheProvider.GetUnmarshalled(ctx, requestUrl, &repository); err == nil {
					return &repository, nil
				}

				if githubToken != "" {
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", githubToken))
				}

				timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Second*30)
				defer cancelTimeout()
				req = req.WithContext(timeoutCtx)

				var rawLanguages external.Languages
				err = httpProvider.ReqUnmarshalledBody(req, &rawLanguages)

				languages := make(model.Language)
				for k, v := range rawLanguages {
					languages[k] = model.LanguageStats{
						Bytes: v,
					}
				}

				repository = model.Repository{
					FullName:      rawRepository.FullName,
					Owner:         rawRepository.Owner.Login,
					Repository:    rawRepository.Name,
					Description:   rawRepository.Description,
					RepositoryUrl: fmt.Sprintf("https://github.com/%s", rawRepository.FullName),
					Languages:     languages,
					License: util.NullableJsonField[string]{
						Value:  rawRepository.License.Value.Key,
						IsNull: rawRepository.License.IsNull,
					},
					Size:      rawRepository.Size,
					CreatedAt: rawRepository.CreatedAt,
					UpdatedAt: rawRepository.UpdatedAt,
				}

				if err == nil {
					_ = cacheProvider.SetMarshalled(ctx, requestUrl, repository, time.Minute*cacheDurationInMin)
				}

				return &repository, nil
			},
		}, nil
	default:
		return nil, errors.New("unsupported github api version")
	}
}
