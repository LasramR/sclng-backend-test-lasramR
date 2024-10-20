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
	githubToken  string
	HttpProvider providers.HttpProvider
}

func (ghRepository *GithubApiRepositoryImpl) GetProjects(ctx context.Context, grb builder.GithubRequestBuilder) (external.RepositoriesResponse, error) {
	grb.Authorization(ghRepository.githubToken)
	req, err := grb.Build(ctx, http.MethodGet, "/search/repositories")
	if err != nil {
		return external.RepositoriesResponse{}, nil
	}

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Second*30)
	defer cancelTimeout()

	req = req.WithContext(timeoutCtx)

	var projects external.RepositoriesResponse
	err = ghRepository.HttpProvider.ReqUnmarshalledBody(req, &projects)

	return projects, err
}

func (ghRepository *GithubApiRepositoryImpl) GetLanguages(ctx context.Context, project external.RepositoriesResponseItem) (external.Languages, error) {
	req, err := http.NewRequest(http.MethodGet, project.LanguagesUrl, nil)
	if err != nil {
		return external.Languages{}, nil
	}

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Second*30)
	defer cancelTimeout()

	if ghRepository.githubToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ghRepository.githubToken))
	}
	req = req.WithContext(timeoutCtx)

	var languages external.Languages
	err = ghRepository.HttpProvider.ReqUnmarshalledBody(req, &languages)

	return languages, err
}

func NewGithubApiRepositoryImpl(githubApiBaseUrl, githubToken string, httpProvider providers.HttpProvider) *GithubApiRepositoryImpl {
	return &GithubApiRepositoryImpl{
		githubToken:  githubToken,
		HttpProvider: httpProvider,
	}
}
