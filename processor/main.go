package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cicd-pipeline-builder/processor/executor"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Initialize Redis
	redisClient := initRedis()
	defer redisClient.Close()

	// Initialize Docker executor
	dockerExecutor := executor.NewDockerExecutor()

	log.Println("Job Processor started, waiting for jobs...")

	// Start processing jobs
	processJobs(redisClient, dockerExecutor)
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

func processJobs(redisClient *redis.Client, dockerExecutor *executor.DockerExecutor) {
	ctx := context.Background()

	for {
		// Try to get a job from the queue
		result, err := redisClient.BRPop(ctx, 0, "job:queue").Result()
		if err != nil {
			log.Printf("Error reading from queue: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if len(result) < 2 {
			continue
		}

		jobID := result[1]
		log.Printf("Processing job: %s", jobID)

		// Execute the job
		if err := dockerExecutor.ExecuteJob(ctx, jobID); err != nil {
			log.Printf("Error executing job %s: %v", jobID, err)
			// Publish error to Redis
			redisClient.Publish(ctx, fmt.Sprintf("job:%s:error", jobID), err.Error())
		}

		log.Printf("Job %s completed", jobID)
	}
}
