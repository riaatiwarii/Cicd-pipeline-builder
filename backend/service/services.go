package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"time"

	"cicd-pipeline-builder/backend/model"
	"cicd-pipeline-builder/backend/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication logic
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo}
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !checkPassword(user.Password, password) {
		return "", errors.New("invalid password")
	}

	return s.generateToken(user)
}

// Register creates a new user
func (s *AuthService) Register(username, email, password, fullName string) (*model.User, error) {
	// Check if user exists
	_, err := s.userRepo.GetByUsername(username)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       uuid.New(),
		Username: username,
		Email:    email,
		Password: hashedPassword,
		FullName: fullName,
		Role:     "user",
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// generateToken creates a JWT token for a user
func (s *AuthService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID.String(),
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// PipelineService handles pipeline business logic
type PipelineService struct {
	pipelineRepo *repository.PipelineRepository
	buildRepo    *repository.BuildRepository
	jobRepo      *repository.JobRepository
	redis        interface{} // Redis client
}

// NewPipelineService creates a new pipeline service
func NewPipelineService(pipelineRepo *repository.PipelineRepository, buildRepo *repository.BuildRepository, jobRepo *repository.JobRepository, redis interface{}) *PipelineService {
	return &PipelineService{pipelineRepo, buildRepo, jobRepo, redis}
}

// ListPipelines retrieves all pipelines for a user
func (s *PipelineService) ListPipelines(userID uuid.UUID) ([]model.Pipeline, error) {
	return s.pipelineRepo.ListByUserID(userID)
}

// GetPipeline retrieves a pipeline by ID
func (s *PipelineService) GetPipeline(id uuid.UUID) (*model.Pipeline, error) {
	return s.pipelineRepo.GetByID(id)
}

// CreatePipeline creates a new pipeline
func (s *PipelineService) CreatePipeline(pipeline *model.Pipeline) error {
	pipeline.ID = uuid.New()
	return s.pipelineRepo.Create(pipeline)
}

// UpdatePipeline updates a pipeline
func (s *PipelineService) UpdatePipeline(pipeline *model.Pipeline) error {
	return s.pipelineRepo.Update(pipeline)
}

// DeletePipeline deletes a pipeline
func (s *PipelineService) DeletePipeline(id uuid.UUID) error {
	return s.pipelineRepo.Delete(id)
}

// TriggerPipeline manually triggers a pipeline
func (s *PipelineService) TriggerPipeline(pipelineID uuid.UUID, userID uuid.UUID) (*model.Build, error) {
	pipeline, err := s.GetPipeline(pipelineID)
	if err != nil {
		return nil, err
	}

	build := &model.Build{
		ID:          uuid.New(),
		PipelineID:  pipelineID,
		Status:      "pending",
		TriggeredBy: &userID,
		TriggerType: "manual",
		CreatedAt:   time.Now(),
	}

	if err := s.buildRepo.Create(build); err != nil {
		return nil, err
	}

	// Queue the job
	// TODO: Send to Redis queue for processor

	return build, nil
}

// BuildService handles build business logic
type BuildService struct {
	buildRepo    *repository.BuildRepository
	jobRepo      *repository.JobRepository
	logRepo      *repository.LogRepository
	pipelineRepo *repository.PipelineRepository
	redis        interface{}
}

// NewBuildService creates a new build service
func NewBuildService(buildRepo *repository.BuildRepository, jobRepo *repository.JobRepository, logRepo *repository.LogRepository, pipelineRepo *repository.PipelineRepository, redis interface{}) *BuildService {
	return &BuildService{buildRepo, jobRepo, logRepo, pipelineRepo, redis}
}

// ListBuilds retrieves all builds for a pipeline
func (s *BuildService) ListBuilds(pipelineID uuid.UUID, limit, offset int) ([]model.Build, error) {
	return s.buildRepo.ListByPipelineID(pipelineID, limit, offset)
}

// GetBuild retrieves a build by ID
func (s *BuildService) GetBuild(id uuid.UUID) (*model.Build, error) {
	return s.buildRepo.GetByID(id)
}

// GetBuildLogs retrieves logs for a build
func (s *BuildService) GetBuildLogs(buildID uuid.UUID) ([]model.Log, error) {
	// Get all jobs for the build
	jobs, err := s.jobRepo.ListByBuildID(buildID)
	if err != nil {
		return nil, err
	}

	var allLogs []model.Log
	for _, job := range jobs {
		logs, err := s.logRepo.ListByJobID(job.ID, 1000)
		if err == nil {
			allLogs = append(allLogs, logs...)
		}
	}

	return allLogs, nil
}

// WebhookService handles webhook business logic
type WebhookService struct {
	webhookRepo      *repository.WebhookRepository
	pipelineService  *PipelineService
}

// NewWebhookService creates a new webhook service
func NewWebhookService(webhookRepo *repository.WebhookRepository, pipelineService *PipelineService) *WebhookService {
	return &WebhookService{webhookRepo, pipelineService}
}

// ListWebhooks retrieves all webhooks for a pipeline
func (s *WebhookService) ListWebhooks(pipelineID uuid.UUID) ([]model.Webhook, error) {
	return s.webhookRepo.ListByPipelineID(pipelineID)
}

// CreateWebhook creates a new webhook
func (s *WebhookService) CreateWebhook(webhook *model.Webhook) error {
	webhook.ID = uuid.New()
	webhook.SecretKey = generateSecretKey()
	return s.webhookRepo.Create(webhook)
}

// TriggerService handles trigger business logic
type TriggerService struct {
	triggerRepo      *repository.TriggerRepository
	pipelineService  *PipelineService
}

// NewTriggerService creates a new trigger service
func NewTriggerService(triggerRepo *repository.TriggerRepository, pipelineService *PipelineService) *TriggerService {
	return &TriggerService{triggerRepo, pipelineService}
}

// ListTriggers retrieves all triggers for a pipeline
func (s *TriggerService) ListTriggers(pipelineID uuid.UUID) ([]model.Trigger, error) {
	return s.triggerRepo.ListByPipelineID(pipelineID)
}

// CreateTrigger creates a new trigger
func (s *TriggerService) CreateTrigger(trigger *model.Trigger) error {
	trigger.ID = uuid.New()
	return s.triggerRepo.Create(trigger)
}

// Helper functions
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateSecretKey() string {
	key := make([]byte, 32)
	rand.Read(key)
	return base64.StdEncoding.EncodeToString(key)
}
