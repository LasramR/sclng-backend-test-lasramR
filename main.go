package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/LasramR/sclng-backend-test-lasramR/api"
	"github.com/LasramR/sclng-backend-test-lasramR/builder"
	"github.com/LasramR/sclng-backend-test-lasramR/providers"
	"github.com/LasramR/sclng-backend-test-lasramR/repositories"
	"github.com/LasramR/sclng-backend-test-lasramR/services"
	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
	redis "github.com/redis/go-redis/v9"
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
	httpProvider := providers.NewNativeHttpProvider(providers.NativeHttpClient{
		Do: http.DefaultClient.Do,
	})

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("redis:%d", cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0, // Use default DB
		Protocol: 2, // Connection protocol
	})

	if rdb.Ping(context.Background()).Err() != nil {
		log.Fatalf("could not connect to redis")
	}

	cacheProvider := providers.NewRedisCacheProvider(&providers.RedisClient{
		Get: rdb.Get,
		Set: rdb.Set,
	})

	githubApiRepository, err := repositories.NewGithubApiRepository(
		builder.GITHUB_API_2022_11_28,
		httpProvider,
		cacheProvider,
		cfg.GithubToken,
	)
	if err != nil {
		log.Fatalf("could not initialize github repository")
	}

	githubService := services.NewGithubServiceImpl(githubApiRepository)

	log.Info("Initializing routes")
	router := handlers.NewRouter(log)
	router.HandleFunc("/repositories", handlers.HandlerFunc(api.GitHubProjectsHandler(githubService, cacheProvider, builder.GITHUB_API_2022_11_28)))
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
