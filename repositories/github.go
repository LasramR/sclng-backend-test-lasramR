package repositories

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/LasramR/sclng-backend-test-lasramR/model"
	"github.com/LasramR/sclng-backend-test-lasramR/providers/http_provider"
)

type GithubApiRepository interface {
	GetProjects(ctx context.Context) (model.GithubApiProjectsResponse, error)
	GetLanguages(ctx context.Context, project model.GithubApiProjectsResponseItem) (model.GithubApiLanguageURLResponse, error)
}

type GithubApiRepositoryImpl struct {
	GithubApiBaseUrl string
	githubToken      string
	HttpProvider     http_provider.HttpProvider
}

func (ghRepository *GithubApiRepositoryImpl) GetProjects(ctx context.Context) (model.GithubApiProjectsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, ghRepository.GithubApiBaseUrl, nil)
	if err != nil {
		return model.GithubApiProjectsResponse{}, nil
	}

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Second*30)
	defer cancelTimeout()

	if ghRepository.githubToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ghRepository.githubToken))
	}
	req = req.WithContext(timeoutCtx)

	var projects model.GithubApiProjectsResponse
	err = ghRepository.HttpProvider.ReqUnmarshalledBody(req, &projects)

	return projects, err
}

func (ghRepository *GithubApiRepositoryImpl) GetLanguages(ctx context.Context, project model.GithubApiProjectsResponseItem) (model.GithubApiLanguageURLResponse, error) {
	req, err := http.NewRequest(http.MethodGet, project.LanguagesUrl, nil)
	if err != nil {
		return model.GithubApiLanguageURLResponse{}, nil
	}

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Second*30)
	defer cancelTimeout()

	if ghRepository.githubToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ghRepository.githubToken))
	}
	req = req.WithContext(timeoutCtx)

	var languages model.GithubApiLanguageURLResponse
	err = ghRepository.HttpProvider.ReqUnmarshalledBody(req, &languages)

	return languages, err
}

func NewGithubApiRepositoryImpl(githubApiBaseUrl, githubToken string, httpProvider http_provider.HttpProvider) *GithubApiRepositoryImpl {
	return &GithubApiRepositoryImpl{
		GithubApiBaseUrl: githubApiBaseUrl,
		githubToken:      githubToken,
		HttpProvider:     httpProvider,
	}
}
