package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/LasramR/sclng-backend-test-lasramR/providers/http_provider"
	"github.com/LasramR/sclng-backend-test-lasramR/repositories"
	"github.com/LasramR/sclng-backend-test-lasramR/services"
	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
)

func test_service() {
	httpProvider := http_provider.NewNativeHttpProvider(&http_provider.NativeHttpClient{
		Do:         http.DefaultClient.Do,
		NewRequest: http.NewRequest,
	})
	githubApiRepository := repositories.NewGithubApiRepositoryImpl(
		"https://api.github.com/search/repositories?q=is:public&per_page=5",
		httpProvider,
	)
	githubService := services.NewGithubServiceImpl(githubApiRepository)

	ctx := context.Background()
	projects, err := githubService.GetGithubProjects(ctx)

	if err != nil {
		fmt.Println("err: ", err)
	} else {
		fmt.Println("projects: ", projects[0])
	}
}

func main() {
	log := logger.Default()
	log.Info("Initializing app")
	cfg, err := newConfig()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration")
		os.Exit(1)
	}

	log.Info("Initializing routes")
	router := handlers.NewRouter(log)
	router.HandleFunc("/ping", pongHandler)

	// Initialize web server and configure the following routes:
	// GET /repos
	// GET /stats

	log.Info("Testing gh service...")
	test_service()

	log = log.WithField("port", cfg.Port)
	log.Info("Listening...")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)

	if err != nil {
		log.WithError(err).Error("Fail to listen to the given port")
		os.Exit(2)
	}
}

func pongHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(map[string]string{"status": "pong"})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}
