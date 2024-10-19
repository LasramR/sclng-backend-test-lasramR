package services

import (
	"context"
	"sync"
	"time"

	"github.com/LasramR/sclng-backend-test-lasramR/model"
	"github.com/LasramR/sclng-backend-test-lasramR/repositories"
	"github.com/LasramR/sclng-backend-test-lasramR/util"
)

type GitHubProject struct {
	FullName      string                                 `json:"full_name"`
	Owner         string                                 `json:"owner"`
	Repository    string                                 `json:"repository"`
	RepositoryUrl string                                 `json:"repository_url"`
	Languages     map[string]GitHubProjectLanguagesStats `json:"languages"`
	License       util.NullableJsonField[string]         `json:"license"`
	Size          int                                    `json:"size"`
	UpdatedAt     string                                 `json:"updated_at"`
}

type GitHubProjectLanguagesStats struct {
	Bytes int `json:"bytes"`
}

type GithubService interface {
	GetGithubProjectsWithStats(ctx context.Context) ([]*GitHubProject, error)
}

type GithubServiceImpl struct {
	GithubRepository repositories.GithubApiRepository
}

func (ghService *GithubServiceImpl) GetGithubProjectsWithStats(ctx context.Context) ([]*GitHubProject, error) {
	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Second*30)
	defer cancelTimeout()

	projectsFromApi, err := ghService.GithubRepository.GetProjects(timeoutCtx)

	if err != nil {
		return nil, err
	}

	projectCount := len(projectsFromApi.Items)
	projects := make([]*GitHubProject, projectCount)
	wg := sync.WaitGroup{}
	responseCh := make(chan *util.IndexedResult[*GitHubProject], projectCount)

	for i, v := range projectsFromApi.Items {
		wg.Add(1)
		go func(rawProject model.GithubApiProjectsResponseItem, idx int) {
			defer wg.Done()

			if err != nil {
				responseCh <- &util.IndexedResult[*GitHubProject]{
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

				languages := make(map[string]GitHubProjectLanguagesStats)
				for k, v := range languagesApiResult {
					languages[k] = GitHubProjectLanguagesStats{
						Bytes: v,
					}
				}
				responseCh <- &util.IndexedResult[*GitHubProject]{
					Value: &GitHubProject{
						FullName:      rawProject.FullName,
						Owner:         rawProject.Owner.Login,
						Repository:    rawProject.Name,
						RepositoryUrl: rawProject.Url,
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
