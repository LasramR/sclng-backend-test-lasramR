package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/LasramR/sclng-backend-test-lasramR/builder"
	"github.com/LasramR/sclng-backend-test-lasramR/model"
	"github.com/LasramR/sclng-backend-test-lasramR/model/external"
	"github.com/LasramR/sclng-backend-test-lasramR/repositories"
	"github.com/LasramR/sclng-backend-test-lasramR/util"
)

type GithubService interface {
	GetGithubProjectsWithStats(ctx context.Context, grb builder.GithubRequestBuilder) ([]*model.Repository, error)
}

type GithubServiceImpl struct {
	GithubRepository repositories.GithubApiRepository
}

func (ghService *GithubServiceImpl) GetGithubProjectsWithStats(ctx context.Context, grb builder.GithubRequestBuilder) ([]*model.Repository, error) {
	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Second*30)
	defer cancelTimeout()

	projectsFromApi, err := ghService.GithubRepository.GetProjects(timeoutCtx, grb)

	if err != nil {
		return nil, err
	}

	projectCount := len(projectsFromApi.Items)
	projects := make([]*model.Repository, projectCount)
	wg := sync.WaitGroup{}
	responseCh := make(chan *util.IndexedResult[*model.Repository], projectCount)

	for i, v := range projectsFromApi.Items {
		wg.Add(1)
		go func(rawProject external.RepositoriesResponseItem, idx int) {
			defer wg.Done()

			if err != nil {
				responseCh <- &util.IndexedResult[*model.Repository]{
					Value: nil,
					Error: err,
					Index: i,
				}
			} else {
				timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Second*30)
				defer cancelTimeout()

				languagesApiResult, err := ghService.GithubRepository.GetLanguages(timeoutCtx, rawProject)

				if err != nil {
				}

				languages := make(model.Language)
				for k, v := range languagesApiResult {
					languages[k] = model.LanguageStats{
						Bytes: v,
					}
				}
				responseCh <- &util.IndexedResult[*model.Repository]{
					Value: &model.Repository{
						FullName:      rawProject.FullName,
						Owner:         rawProject.Owner.Login,
						Repository:    rawProject.Name,
						RepositoryUrl: fmt.Sprintf("https://github.com/%s", rawProject.FullName),
						Languages:     languages,
						License: util.NullableJsonField[string]{
							Value:  rawProject.License.Value.Key,
							IsNull: rawProject.License.IsNull,
						},
						Size:      rawProject.Size,
						UpdatedAt: rawProject.UpdatedAt,
					},
					Error: nil,
					Index: i,
				}
			}
		}(v, i)
	}

	go func() {
		wg.Wait()
		close(responseCh)
	}()

	for result := range responseCh {
		if result.Error == nil {
			projects[result.Index] = result.Value

		}
		// TODO handle error ?
	}

	return projects, nil
}

func NewGithubServiceImpl(ghRepository repositories.GithubApiRepository) *GithubServiceImpl {
	return &GithubServiceImpl{
		GithubRepository: ghRepository,
	}
}
