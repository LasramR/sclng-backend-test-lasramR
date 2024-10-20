package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/LasramR/sclng-backend-test-lasramR/builder"
	"github.com/LasramR/sclng-backend-test-lasramR/model"
	"github.com/LasramR/sclng-backend-test-lasramR/model/version"
	"github.com/LasramR/sclng-backend-test-lasramR/providers"
	"github.com/LasramR/sclng-backend-test-lasramR/repositories"
	"github.com/LasramR/sclng-backend-test-lasramR/services"
	"github.com/LasramR/sclng-backend-test-lasramR/util"
	"github.com/Scalingo/go-utils/logger"
)

// Compute error object and marshal it in request response writer
func errorFallback(w http.ResponseWriter, errs []string, status int) error {
	response := model.ApiError{
		Status: status,
		Reason: errs,
	}
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(response)
}

// Compute success object and marshal it in request response writer
func successFallback(w http.ResponseWriter, r *http.Request, repos repositories.GithubRepositoriesResult) error {
	response := model.ApiListResponse[[]*model.Repository]{
		TotalCount:       repos.Total,
		Count:            len(repos.Repositories),
		Content:          repos.Repositories,
		IncompleteResult: repos.IncompleteResult,
		Page:             0,
		Previous: util.NullableJsonField[string]{
			IsNull: r.URL.Query().Get("page") == "",
			Value:  util.PreviousFullUrlFromRequest(r),
		},
		Next: util.NextFullUrlFromRequest(r),
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(response)
}

// /repos HTTP handle
func GitHubProjectsHandler(
	githubService services.GithubService,
	cacheProvider providers.CacheProvider,
	apiVersion version.GithubAPIVersion,
) util.ScalingoHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
		ctx := r.Context()
		log := logger.Get(ctx)

		// Only respond to GET
		if r.Method != http.MethodGet {
			return errorFallback(w, []string{"GET only endpoint"}, http.StatusMethodNotAllowed)
		}

		requestUrl := util.FullUrlFromRequest(r)
		var repos repositories.GithubRepositoriesResult = repositories.GithubRepositoriesResult{}
		// Returns if successful cache read from requestUrl
		if err := cacheProvider.GetUnmarshalled(ctx, requestUrl, &repos); err == nil {
			return successFallback(w, r, repos)
		}

		grb, err := builder.NewGithubRequestBuilder(apiVersion)

		if err != nil {
			log.WithError(err).Error(err)
			return errorFallback(w, []string{err.Error()}, http.StatusServiceUnavailable)
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
				return errorFallback(w, []string{err.Error()}, http.StatusBadRequest)
			}

			if err := grb.Limit(parsedLimit); err != nil {
				log.WithError(err).Error(err)
				return errorFallback(w, []string{err.Error()}, http.StatusBadRequest)
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
				return errorFallback(w, []string{err.Error()}, http.StatusBadRequest)
			}

			if err := grb.Page(parsedPage); err != nil {
				log.WithError(err).Error(err)
				return errorFallback(w, []string{err.Error()}, http.StatusBadRequest)
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

		// GIVE ME THESE REPOSITORIES
		repos, err = githubService.GetGithubProjectsWithStats(ctx, grb)

		if err != nil {
			log.WithError(err).Error(err)
			return errorFallback(w, []string{err.Error()}, http.StatusInternalServerError)
		}

		// Set in cache
		_ = cacheProvider.SetMarshalled(ctx, requestUrl, repos, time.Minute*5)

		return successFallback(w, r, repos)
	}

}
