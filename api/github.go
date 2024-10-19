package api

import (
	"encoding/json"
	"net/http"

	"github.com/LasramR/sclng-backend-test-lasramR/services"
	"github.com/Scalingo/go-utils/logger"
)

type ScalingoHandlerFunc func(w http.ResponseWriter, r *http.Request, _ map[string]string) error

func GitHubProjectsHandler(
	githubService services.GithubService,
) ScalingoHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
		ctx := r.Context()
		log := logger.Get(ctx)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		projects, err := githubService.GetGithubProjectsWithStats(ctx)

		if err != nil {
			log.WithError(err).Error("Failed to retrieve projects")
		}

		return json.NewEncoder(w).Encode(projects)
	}

}
