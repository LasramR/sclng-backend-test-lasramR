package services

import (
	"context"
	"time"

	"github.com/LasramR/sclng-backend-test-lasramR/builder"
	"github.com/LasramR/sclng-backend-test-lasramR/repositories"
)

type GithubService interface {
	GetGithubProjectsWithStats(ctx context.Context, grb builder.GithubRequestBuilder) (repositories.GithubRepositoriesResult, error)
}

type GithubServiceImpl struct {
	GithubRepository repositories.GithubApiRepository
}

func (ghService *GithubServiceImpl) GetGithubProjectsWithStats(ctx context.Context, grb builder.GithubRequestBuilder) (repositories.GithubRepositoriesResult, error) {
	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Second*30)
	defer cancelTimeout()

	result, err := ghService.GithubRepository.GetManyRepositories(timeoutCtx, grb)

	if err != nil {
		return repositories.GithubRepositoriesResult{}, err
	}

	return result, nil
}

func NewGithubServiceImpl(ghRepository repositories.GithubApiRepository) *GithubServiceImpl {
	return &GithubServiceImpl{
		GithubRepository: ghRepository,
	}
}
