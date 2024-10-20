package repositories

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/LasramR/sclng-backend-test-lasramR/builder"
	"github.com/LasramR/sclng-backend-test-lasramR/model/external"
	"github.com/LasramR/sclng-backend-test-lasramR/providers"
)

type GithubApiRepository interface {
	GetProjects(ctx context.Context, grb builder.GithubRequestBuilder) (external.RepositoriesResponse, error)
	GetLanguages(ctx context.Context, project external.RepositoriesResponseItem) (external.Languages, error)
}

type GithubApiRepositoryImpl struct {
	githubToken   string
	HttpProvider  providers.HttpProvider
	CacheProvider providers.CacheProvider
}

func (ghRepository *GithubApiRepositoryImpl) GetProjects(ctx context.Context, grb builder.GithubRequestBuilder) (external.RepositoriesResponse, error) {
	grb.Authorization(ghRepository.githubToken)
	req, err := grb.Build(ctx, http.MethodGet, "/search/repositories")

	if err != nil {
		return external.RepositoriesResponse{}, nil
	}

	requestUrl := req.URL.String()
	var projects external.RepositoriesResponse

	if err = ghRepository.CacheProvider.GetUnmarshalled(ctx, requestUrl, &projects); err == nil {
		return projects, nil
	}

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Second*30)
	defer cancelTimeout()

	req = req.WithContext(timeoutCtx)
	err = ghRepository.HttpProvider.ReqUnmarshalledBody(req, &projects)

	if err == nil {
		_ = ghRepository.CacheProvider.SetMarshalled(ctx, requestUrl, projects, time.Minute*5)
	}

	return projects, err
}

func (ghRepository *GithubApiRepositoryImpl) GetLanguages(ctx context.Context, project external.RepositoriesResponseItem) (external.Languages, error) {
	req, err := http.NewRequest(http.MethodGet, project.LanguagesUrl, nil)
	if err != nil {
		return external.Languages{}, nil
	}

	requestUrl := project.LanguagesUrl
	var languages external.Languages

	if err = ghRepository.CacheProvider.GetUnmarshalled(ctx, requestUrl, &languages); err == nil {
		return languages, nil
	}

	if ghRepository.githubToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ghRepository.githubToken))
	}

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Second*30)
	defer cancelTimeout()
	req = req.WithContext(timeoutCtx)

	err = ghRepository.HttpProvider.ReqUnmarshalledBody(req, &languages)

	if err == nil {
		_ = ghRepository.CacheProvider.SetMarshalled(ctx, requestUrl, languages, time.Minute*5)
	}

	return languages, err
}

func NewGithubApiRepositoryImpl(githubApiBaseUrl, githubToken string, httpProvider providers.HttpProvider, cacheProvider providers.CacheProvider) *GithubApiRepositoryImpl {
	return &GithubApiRepositoryImpl{
		githubToken:   githubToken,
		HttpProvider:  httpProvider,
		CacheProvider: cacheProvider,
	}
}
