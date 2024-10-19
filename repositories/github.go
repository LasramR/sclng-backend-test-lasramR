package repositories

import (
	"context"

	"github.com/LasramR/sclng-backend-test-lasramR/model"
	"github.com/LasramR/sclng-backend-test-lasramR/providers/http_provider"
)

type GithubApiRepository interface {
	GetProjects(ctx context.Context) (model.GithubApiProjectsResponse, error)
	GetLanguages(ctx context.Context, project model.GithubApiProjectsResponseItem) (model.GithubApiLanguageURLResponse, error)
}

type GithubApiRepositoryImpl struct {
	GithubApiBaseUrl string
	HttpProvider     http_provider.HttpProvider
}

func (ghRepository *GithubApiRepositoryImpl) GetProjects(ctx context.Context) (model.GithubApiProjectsResponse, error) {
	var projects model.GithubApiProjectsResponse
	err := ghRepository.HttpProvider.GetJson(ctx, ghRepository.GithubApiBaseUrl, &projects)
	return projects, err
}

func (ghRepository *GithubApiRepositoryImpl) GetLanguages(ctx context.Context, project model.GithubApiProjectsResponseItem) (model.GithubApiLanguageURLResponse, error) {
	var languages model.GithubApiLanguageURLResponse
	err := ghRepository.HttpProvider.GetJson(ctx, project.LanguagesUrl, &languages)
	return languages, err
}

func NewGithubApiRepositoryImpl(githubApiBaseUrl string, httpProvider http_provider.HttpProvider) *GithubApiRepositoryImpl {
	return &GithubApiRepositoryImpl{
		GithubApiBaseUrl: githubApiBaseUrl,
		HttpProvider:     httpProvider,
	}
}
