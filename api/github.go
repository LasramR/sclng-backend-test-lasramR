package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/LasramR/sclng-backend-test-lasramR/builder"
	"github.com/LasramR/sclng-backend-test-lasramR/services"
	"github.com/LasramR/sclng-backend-test-lasramR/util"
	"github.com/Scalingo/go-utils/logger"
)

type BadRequestBody struct {
	Errors []string `json:"errors"`
}

func GitHubProjectsHandler(
	githubService services.GithubService,
	apiVersion builder.GithubAPIVersion,
) util.ScalingoHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
		ctx := r.Context()
		log := logger.Get(ctx)

		grb, err := builder.NewGithubRequestBuilder(apiVersion)

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			log.WithError(err).Error("Unsupported Github API")
			return err
		}

		queryParams := r.URL.Query()
		queryErrors := BadRequestBody{
			Errors: make([]string, 0, len(queryParams)),
		}
		for k, v := range queryParams {
			err := grb.With(k, strings.Join(v, " "))

			if err != nil {
				log.WithError(err).Error("bad request")
				queryErrors.Errors = append(queryErrors.Errors, err.Error())
			}
		}

		if len(queryErrors.Errors) != 0 {
			// TODO refactor this
			w.WriteHeader(http.StatusBadRequest)
			log.WithError(err).Error("bad request")
			return json.NewEncoder(w).Encode(queryErrors)
		}

		projects, err := githubService.GetGithubProjectsWithStats(ctx, grb)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.WithError(err).Error("Failed to retrieve projects")
			return err
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		return json.NewEncoder(w).Encode(projects)
	}

}
