package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/LasramR/sclng-backend-test-lasramR/api"
	"github.com/LasramR/sclng-backend-test-lasramR/providers/http_provider"
	"github.com/LasramR/sclng-backend-test-lasramR/repositories"
	"github.com/LasramR/sclng-backend-test-lasramR/services"
	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
)

func main() {
	log := logger.Default()

	log.Info("Initializing app")
	cfg, err := newConfig()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration")
		os.Exit(1)
	}

	if cfg.GithubToken == "" {
		log.Warn("Booting without the use of a Github token: the application will run in limited mode")
	}

	log.Info("Initializing services")
	httpProvider := http_provider.NewNativeHttpProvider(http_provider.NativeHttpClient{
		Do: http.DefaultClient.Do,
	})
	githubApiRepository := repositories.NewGithubApiRepositoryImpl(
		"https://api.github.com/search/repositories?q=is:public&per_page=50",
		cfg.GithubToken,
		httpProvider,
	)
	githubService := services.NewGithubServiceImpl(githubApiRepository)

	log.Info("Initializing routes")
	router := handlers.NewRouter(log)
	router.HandleFunc("/projects", handlers.HandlerFunc(api.GitHubProjectsHandler(githubService)))
	// GET /repos
	// GET /stats

	log = log.WithField("port", cfg.Port)
	log.Info("Listening...")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)

	if err != nil {
		log.WithError(err).Error("Fail to listen to the given port")
		os.Exit(2)
	}
}
