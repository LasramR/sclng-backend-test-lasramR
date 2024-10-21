package services

import (
	"context"
	"testing"

	"github.com/LasramR/sclng-backend-test-lasramR/builder"
	"github.com/LasramR/sclng-backend-test-lasramR/model"
	"github.com/LasramR/sclng-backend-test-lasramR/model/version"
	"github.com/LasramR/sclng-backend-test-lasramR/repositories"
	"github.com/LasramR/sclng-backend-test-lasramR/util"
)

type MockGithubRepository struct {
	err error
}

func (mgr *MockGithubRepository) GetManyRepositories(_ctx context.Context, _grb builder.GithubRequestBuilder) (repositories.GithubRepositoriesResult, error) {
	if mgr.err != nil {
		return repositories.GithubRepositoriesResult{}, mgr.err
	}

	return repositories.GithubRepositoriesResult{
		Total:            1,
		IncompleteResult: false,
		Repositories: []*model.Repository{
			{
				FullName:    "fmuiin14/BlazingTool",
				Owner:       "fmuiin14",
				Description: "Brute force ethereum wallet mnemonics",
				Repository:  "fmuiin14/BlazingTool",
				License: util.NullableJsonField[string]{
					Value:  "mit",
					IsNull: false,
				},
				RepositoryUrl: "https://github.com/fmuiin14/BlazingTool",
				CreatedAt:     "2024-10-19T10:17:16Z",
				UpdatedAt:     "2024-10-20T16:36:13Z",
				Languages: model.Language{
					"JavaScript": model.LanguageStats{
						Bytes: 1548,
					},
					"SCSS": model.LanguageStats{
						Bytes: 250,
					},
				},
				Size: 156464,
			},
		},
	}, nil
}

func TestGetGithubProjectsWithStats(t *testing.T) {
	gs := NewGithubService(&MockGithubRepository{})
	grb, _ := builder.NewGithubRequestBuilder(version.GITHUB_API_2022_11_28)
	result, err := gs.GetGithubProjectsWithStats(context.Background(), grb)

	if err != nil {
		t.Fatalf("Should not have returned an error")
	}

	if len(result.Repositories) != 1 || result.Repositories[0].Repository != "fmuiin14/BlazingTool" {
		t.Fatalf("Should have returned Github repository GetManyRepositories result")
	}
}
