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
	"github.com/LasramR/sclng-backend-test-lasramR/providers"
	"github.com/LasramR/sclng-backend-test-lasramR/util"
)

type GithubRepositoriesResult struct {
	Repositories     []*model.Repository `json:"repositories"`
	Total            int                 `json:"total"`
	IncompleteResult bool                `json:"incomplete_result"`
}

type GithubApiRepository interface {
	GetManyRepositories(ctx context.Context, grb builder.GithubRequestBuilder) (GithubRepositoriesResult, error)
}

type githubVersionnedApiRepository[T any, M util.Mappable[T]] struct {
	githubToken   string
	httpProvider  providers.HttpProvider
	cacheProvider providers.CacheProvider
	mapperFunc    util.MapperFunc[T, *model.Repository]
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

	if len(errorsCollected) == 0 {
		_ = gr.cacheProvider.SetMarshalled(ctx, requestUrl, repositories, time.Minute*5)
	}

	return GithubRepositoriesResult{
		Repositories:     mapped,
		Total:            apiResponse.Count(),
		IncompleteResult: len(errorsCollected) != 0,
	}, nil
}

func NewGithubApiRepository(apiVersion builder.GithubAPIVersion, httpProvider providers.HttpProvider, cacheProvider providers.CacheProvider, githubToken string) (GithubApiRepository, error) {
	switch apiVersion {
	case builder.GITHUB_API_2022_11_28:
		return &githubVersionnedApiRepository[external.RepositoriesResponseItem, external.RepositoriesResponse]{
			githubToken:   githubToken,
			httpProvider:  httpProvider,
			cacheProvider: cacheProvider,
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
					RepositoryUrl: fmt.Sprintf("https://github.com/%s", rawRepository.FullName),
					Languages:     languages,
					License: util.NullableJsonField[string]{
						Value:  rawRepository.License.Value.Key,
						IsNull: rawRepository.License.IsNull,
					},
					Size:      rawRepository.Size,
					UpdatedAt: rawRepository.UpdatedAt,
				}

				if err == nil {
					_ = cacheProvider.SetMarshalled(ctx, requestUrl, repository, time.Minute*5)
				}

				return &repository, nil
			},
		}, nil
	default:
		return nil, errors.New("unsupported github api version")
	}
}
