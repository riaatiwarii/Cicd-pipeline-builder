package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"cicd-pipeline-builder/backend/handler"
	"cicd-pipeline-builder/backend/middleware"
	"cicd-pipeline-builder/backend/model"
	"cicd-pipeline-builder/backend/repository"
	"cicd-pipeline-builder/backend/service"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Database
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Redis
	redisClient := initRedis()

	// Auto migrate models
	db.AutoMigrate(
		&model.User{},
		&model.Pipeline{},
		&model.Stage{},
		&model.Build{},
		&model.Job{},
		&model.Log{},
		&model.Webhook{},
		&model.Trigger{},
		&model.Artifact{},
		&model.Credential{},
		&model.EnvVariable{},
		&model.WebhookEvent{},
	)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	pipelineRepo := repository.NewPipelineRepository(db)
	buildRepo := repository.NewBuildRepository(db)
	jobRepo := repository.NewJobRepository(db)
	logRepo := repository.NewLogRepository(db)
	webhookRepo := repository.NewWebhookRepository(db)
	triggerRepo := repository.NewTriggerRepository(db)
	artifactRepo := repository.NewArtifactRepository(db)
	credentialRepo := repository.NewCredentialRepository(db)
	envVarRepo := repository.NewEnvVariableRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	pipelineService := service.NewPipelineService(pipelineRepo, buildRepo, jobRepo, redisClient)
	buildService := service.NewBuildService(buildRepo, jobRepo, logRepo, pipelineRepo, redisClient)
	webhookService := service.NewWebhookService(webhookRepo, pipelineService)
	triggerService := service.NewTriggerService(triggerRepo, pipelineService)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	pipelineHandler := handler.NewPipelineHandler(pipelineService, buildService)
	buildHandler := handler.NewBuildHandler(buildService)
	webhookHandler := handler.NewWebhookHandler(webhookService)
	triggerHandler := handler.NewTriggerHandler(triggerService)
	healthHandler := handler.NewHealthHandler(db, redisClient)

	// Setup Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware())

	// Health check
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	// Auth endpoints (no auth required)
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// API endpoints with auth middleware
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		// Pipeline endpoints
		pipelines := api.Group("/pipelines")
		{
			pipelines.GET("", pipelineHandler.ListPipelines)
			pipelines.POST("", pipelineHandler.CreatePipeline)
			pipelines.GET("/:id", pipelineHandler.GetPipeline)
			pipelines.PUT("/:id", pipelineHandler.UpdatePipeline)
			pipelines.DELETE("/:id", pipelineHandler.DeletePipeline)
			pipelines.POST("/:id/trigger", pipelineHandler.TriggerPipeline)
			pipelines.GET("/:id/builds", buildHandler.GetPipelineBuilds)
		}

		// Build endpoints
		builds := api.Group("/builds")
		{
			builds.GET("", buildHandler.ListBuilds)
			builds.GET("/:id", buildHandler.GetBuild)
			builds.GET("/:id/logs", buildHandler.GetBuildLogs)
			builds.GET("/:id/artifacts", buildHandler.GetBuildArtifacts)
			builds.POST("/:id/cancel", buildHandler.CancelBuild)
		}

		// Webhook endpoints
		webhooks := api.Group("/webhooks")
		{
			webhooks.GET("", webhookHandler.ListWebhooks)
			webhooks.POST("", webhookHandler.CreateWebhook)
			webhooks.GET("/:id", webhookHandler.GetWebhook)
			webhooks.PUT("/:id", webhookHandler.UpdateWebhook)
			webhooks.DELETE("/:id", webhookHandler.DeleteWebhook)
		}

		// Trigger endpoints
		triggers := api.Group("/triggers")
		{
			triggers.GET("", triggerHandler.ListTriggers)
			triggers.POST("", triggerHandler.CreateTrigger)
			triggers.GET("/:id", triggerHandler.GetTrigger)
			triggers.PUT("/:id", triggerHandler.UpdateTrigger)
			triggers.DELETE("/:id", triggerHandler.DeleteTrigger)
		}
	}

	// Public webhook endpoints (no auth)
	public := router.Group("/api/v1/public")
	{
		public.POST("/webhooks/github", webhookHandler.HandleGitHubWebhook)
		public.POST("/webhooks/gitlab", webhookHandler.HandleGitLabWebhook)
	}

	log.Printf("Starting server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://cicd:cicd_password_dev@localhost:5432/cicd_db?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func initRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	return redis.NewClient(opts)
}
