package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/LasramR/sclng-backend-test-lasramR/builder"
	"github.com/LasramR/sclng-backend-test-lasramR/model"
	"github.com/LasramR/sclng-backend-test-lasramR/services"
	"github.com/LasramR/sclng-backend-test-lasramR/util"
	"github.com/Scalingo/go-utils/logger"
)

func errorFallback[T any](w http.ResponseWriter, err T, status int) error {
	response := model.ApiError[T]{
		Status: status,
		Reason: err,
	}
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(response)
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
			log.WithError(err).Error(err)
			return errorFallback(w, err.Error(), http.StatusServiceUnavailable)
		}

		queryParams := r.URL.Query()

		// Setting results limit if set in query
		limit := queryParams.Get("limit")
		if limit != "" {
			queryParams.Del("limit")
			parsedLimit, err := strconv.Atoi(limit)
			if err != nil {
				err = errors.New("invalid limit parameter")
				log.WithError(err).Error(err)
				return errorFallback(w, err.Error(), http.StatusBadRequest)
			}

			if err := grb.Limit(parsedLimit); err != nil {
				log.WithError(err).Error(err)
				return errorFallback(w, err.Error(), http.StatusBadRequest)
			}
		}

		// Setting results page if set in query
		page := queryParams.Get("page")
		if page != "" {
			queryParams.Del("page")
			parsedPage, err := strconv.Atoi(page)
			if err != nil {
				err = errors.New("invalid page parameter")
				log.WithError(err).Error(err)
				return errorFallback(w, err.Error(), http.StatusBadRequest)
			}

			if err := grb.Page(parsedPage); err != nil {
				log.WithError(err).Error(err)
				return errorFallback(w, err.Error(), http.StatusBadRequest)
			}
		}

		// Consumming leftovers query parameters
		queryParamsErrors := make([]string, 0, len(queryParams))
		for k, v := range queryParams {
			err := grb.With(k, strings.Join(v, " "))

			if err != nil {
				queryParamsErrors = append(queryParamsErrors, err.Error())
			}
		}

		// If we collected errors
		if len(queryParamsErrors) != 0 {
			log.WithError(err).Error(err)
			return errorFallback(w, queryParamsErrors, http.StatusBadRequest)
		}

		// GIVE ME THESE PROJECTS
		projects, err := githubService.GetGithubProjectsWithStats(ctx, grb)

		if err != nil {
			log.WithError(err).Error(err)
			return errorFallback(w, err.Error(), http.StatusInternalServerError)
		}

		Previous := util.NullableJsonField[string]{
			IsNull: true,
			Value:  "",
		}
		if page != "" {
			Previous.IsNull = false
			Previous.Value = util.PreviousFullUrlFromRequest(r)
		}

		// TODO include metadatas about content total count
		response := model.ApiResponse[[]*model.Repository]{
			Count:            len(projects),
			Content:          projects,
			IncompleteResult: true,
			Page:             0,
			Previous:         Previous,
			Next:             util.NextFullUrlFromRequest(r),
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		return json.NewEncoder(w).Encode(response)
	}

}
