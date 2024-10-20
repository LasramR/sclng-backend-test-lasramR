package services

import (
	"context"
	"time"

	"github.com/LasramR/sclng-backend-test-lasramR/builder"
	"github.com/LasramR/sclng-backend-test-lasramR/repositories"
)

// Github service business logic
type GithubService interface {
	// Returns repositories with computed stats from GithubAPIRepository
	GetGithubProjectsWithStats(ctx context.Context, grb builder.GithubRequestBuilder) (repositories.GithubRepositoriesResult, error)
}

type GithubServiceImpl struct {
	GithubRepository repositories.GithubApiRepository
}

func (gs *GithubServiceImpl) GetGithubProjectsWithStats(ctx context.Context, grb builder.GithubRequestBuilder) (repositories.GithubRepositoriesResult, error) {
	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Second*30)
	defer cancelTimeout()

	result, err := gs.GithubRepository.GetManyRepositories(timeoutCtx, grb)

	if err != nil {
		return repositories.GithubRepositoriesResult{}, err
	}

	return result, nil
}

func NewGithubService(gr repositories.GithubApiRepository) GithubService {
	return &GithubServiceImpl{
		GithubRepository: gr,
	}
}
