package executor

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// DockerExecutor executes jobs in Docker containers
type DockerExecutor struct {
	client *client.Client
}

// NewDockerExecutor creates a new Docker executor
func NewDockerExecutor() *DockerExecutor {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	return &DockerExecutor{client: cli}
}

// ExecuteJob executes a job in a Docker container
func (e *DockerExecutor) ExecuteJob(ctx context.Context, jobID string) error {
	// This would fetch the job details from the database
	// For now, this is a placeholder that shows the execution flow

	log.Printf("Executing job: %s", jobID)

	// TODO: Fetch job details from database
	// job := fetchJobFromDB(jobID)

	// TODO: Create and run container
	// err := e.runContainer(ctx, job)

	return nil
}

// RunContainer runs a Docker container for a stage
func (e *DockerExecutor) RunContainer(ctx context.Context, image string, script string, workDir string) (int, error) {
	// Pull image if needed
	if err := e.pullImage(ctx, image); err != nil {
		return 1, fmt.Errorf("failed to pull image: %w", err)
	}

	// Create container
	resp, err := e.client.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			Cmd:   []string{"/bin/sh", "-c", script},
			WorkingDir: "/workspace",
		},
		&container.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:/workspace", workDir),
			},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		return 1, fmt.Errorf("failed to create container: %w", err)
	}

	containerID := resp.ID
	log.Printf("Created container: %s", containerID)

	// Start container
	if err := e.client.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return 1, fmt.Errorf("failed to start container: %w", err)
	}

	// Stream logs
	logs, err := e.client.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return 1, fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logs.Close()

	// Print logs as they come in
	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		log.Println(scanner.Text())
		// TODO: Send to database/WebSocket for real-time display
	}

	// Wait for container to finish
	statusCh, errCh := e.client.ContainerWait(ctx, containerID, container.WaitConditionNextExit)
	select {
	case status := <-statusCh:
		log.Printf("Container exited with status: %d", status.StatusCode)
		
		// Clean up
		e.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
		
		return int(status.StatusCode), nil
	case err := <-errCh:
		return 1, fmt.Errorf("error waiting for container: %w", err)
	}
}

// pullImage pulls a Docker image
func (e *DockerExecutor) pullImage(ctx context.Context, image string) error {
	log.Printf("Pulling image: %s", image)
	
	reader, err := e.client.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	// Consume the reader
	_, err = io.Copy(os.Stdout, reader)
	return err
}

// ExecuteShellScript executes a shell script locally (for testing)
func ExecuteShellScript(script string, workDir string) (int, error) {
	cmd := exec.Command("/bin/sh", "-c", script)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), nil
		}
		return 1, err
	}

	return 0, nil
}

// PipelineExecutor executes a complete pipeline
type PipelineExecutor struct {
	dockerExecutor *DockerExecutor
}

// NewPipelineExecutor creates a new pipeline executor
func NewPipelineExecutor() *PipelineExecutor {
	return &PipelineExecutor{
		dockerExecutor: NewDockerExecutor(),
	}
}

// ExecutePipeline executes all stages in a pipeline
func (e *PipelineExecutor) ExecutePipeline(ctx context.Context, stages []Stage, workDir string) error {
	for i, stage := range stages {
		log.Printf("Executing stage %d: %s", i+1, stage.Name)

		// Create timeout context
		stageCtx, cancel := context.WithTimeout(ctx, time.Duration(stage.TimeoutSeconds)*time.Second)
		
		exitCode, err := e.dockerExecutor.RunContainer(stageCtx, stage.DockerImage, stage.Script, workDir)
		cancel()

		if err != nil {
			log.Printf("Stage %s failed: %v", stage.Name, err)
			if !stage.AllowFailure {
				return fmt.Errorf("stage %s failed", stage.Name)
			}
		}

		if exitCode != 0 && !stage.AllowFailure {
			return fmt.Errorf("stage %s exited with code %d", stage.Name, exitCode)
		}

		log.Printf("Stage %s completed successfully", stage.Name)
	}

	return nil
}

// Stage represents a pipeline stage
type Stage struct {
	Name          string
	DockerImage   string
	Script        string
	TimeoutSeconds int
	AllowFailure  bool
}
