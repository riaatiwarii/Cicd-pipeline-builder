package repository

import (
	"cicd-pipeline-builder/backend/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update updates a user
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// PipelineRepository handles pipeline database operations
type PipelineRepository struct {
	db *gorm.DB
}

// NewPipelineRepository creates a new pipeline repository
func NewPipelineRepository(db *gorm.DB) *PipelineRepository {
	return &PipelineRepository{db}
}

// ListByUserID retrieves all pipelines for a user
func (r *PipelineRepository) ListByUserID(userID uuid.UUID) ([]model.Pipeline, error) {
	var pipelines []model.Pipeline
	err := r.db.Preload("Stages").Preload("Builds").Where("user_id = ?", userID).Find(&pipelines).Error
	return pipelines, err
}

// GetByID retrieves a pipeline by ID
func (r *PipelineRepository) GetByID(id uuid.UUID) (*model.Pipeline, error) {
	var pipeline model.Pipeline
	err := r.db.Preload("Stages").Preload("Builds").First(&pipeline, "id = ?", id).Error
	return &pipeline, err
}

// Create creates a new pipeline
func (r *PipelineRepository) Create(pipeline *model.Pipeline) error {
	return r.db.Create(pipeline).Error
}

// Update updates a pipeline
func (r *PipelineRepository) Update(pipeline *model.Pipeline) error {
	return r.db.Save(pipeline).Error
}

// Delete deletes a pipeline
func (r *PipelineRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Pipeline{}, "id = ?", id).Error
}

// BuildRepository handles build database operations
type BuildRepository struct {
	db *gorm.DB
}

// NewBuildRepository creates a new build repository
func NewBuildRepository(db *gorm.DB) *BuildRepository {
	return &BuildRepository{db}
}

// ListByPipelineID retrieves all builds for a pipeline
func (r *BuildRepository) ListByPipelineID(pipelineID uuid.UUID, limit, offset int) ([]model.Build, error) {
	var builds []model.Build
	err := r.db.Preload("Jobs").Where("pipeline_id = ?", pipelineID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&builds).Error
	return builds, err
}

// GetByID retrieves a build by ID
func (r *BuildRepository) GetByID(id uuid.UUID) (*model.Build, error) {
	var build model.Build
	err := r.db.Preload("Jobs").Preload("Artifacts").First(&build, "id = ?", id).Error
	return &build, err
}

// Create creates a new build
func (r *BuildRepository) Create(build *model.Build) error {
	return r.db.Create(build).Error
}

// Update updates a build
func (r *BuildRepository) Update(build *model.Build) error {
	return r.db.Save(build).Error
}

// JobRepository handles job database operations
type JobRepository struct {
	db *gorm.DB
}

// NewJobRepository creates a new job repository
func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db}
}

// GetByID retrieves a job by ID
func (r *JobRepository) GetByID(id uuid.UUID) (*model.Job, error) {
	var job model.Job
	err := r.db.First(&job, "id = ?", id).Error
	return &job, err
}

// ListByBuildID retrieves all jobs for a build
func (r *JobRepository) ListByBuildID(buildID uuid.UUID) ([]model.Job, error) {
	var jobs []model.Job
	err := r.db.Where("build_id = ?", buildID).Order("created_at ASC").Find(&jobs).Error
	return jobs, err
}

// Create creates a new job
func (r *JobRepository) Create(job *model.Job) error {
	return r.db.Create(job).Error
}

// Update updates a job
func (r *JobRepository) Update(job *model.Job) error {
	return r.db.Save(job).Error
}

// LogRepository handles log database operations
type LogRepository struct {
	db *gorm.DB
}

// NewLogRepository creates a new log repository
func NewLogRepository(db *gorm.DB) *LogRepository {
	return &LogRepository{db}
}

// ListByJobID retrieves all logs for a job
func (r *LogRepository) ListByJobID(jobID uuid.UUID, limit int) ([]model.Log, error) {
	var logs []model.Log
	err := r.db.Where("job_id = ?", jobID).Order("line_number ASC").Limit(limit).Find(&logs).Error
	return logs, err
}

// Create creates a new log entry
func (r *LogRepository) Create(log *model.Log) error {
	return r.db.Create(log).Error
}

// WebhookRepository handles webhook database operations
type WebhookRepository struct {
	db *gorm.DB
}

// NewWebhookRepository creates a new webhook repository
func NewWebhookRepository(db *gorm.DB) *WebhookRepository {
	return &WebhookRepository{db}
}

// ListByPipelineID retrieves all webhooks for a pipeline
func (r *WebhookRepository) ListByPipelineID(pipelineID uuid.UUID) ([]model.Webhook, error) {
	var webhooks []model.Webhook
	err := r.db.Where("pipeline_id = ?", pipelineID).Find(&webhooks).Error
	return webhooks, err
}

// GetByID retrieves a webhook by ID
func (r *WebhookRepository) GetByID(id uuid.UUID) (*model.Webhook, error) {
	var webhook model.Webhook
	err := r.db.First(&webhook, "id = ?", id).Error
	return &webhook, err
}

// Create creates a new webhook
func (r *WebhookRepository) Create(webhook *model.Webhook) error {
	return r.db.Create(webhook).Error
}

// Update updates a webhook
func (r *WebhookRepository) Update(webhook *model.Webhook) error {
	return r.db.Save(webhook).Error
}

// Delete deletes a webhook
func (r *WebhookRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Webhook{}, "id = ?", id).Error
}

// TriggerRepository handles trigger database operations
type TriggerRepository struct {
	db *gorm.DB
}

// NewTriggerRepository creates a new trigger repository
func NewTriggerRepository(db *gorm.DB) *TriggerRepository {
	return &TriggerRepository{db}
}

// ListByPipelineID retrieves all triggers for a pipeline
func (r *TriggerRepository) ListByPipelineID(pipelineID uuid.UUID) ([]model.Trigger, error) {
	var triggers []model.Trigger
	err := r.db.Where("pipeline_id = ?", pipelineID).Find(&triggers).Error
	return triggers, err
}

// GetByID retrieves a trigger by ID
func (r *TriggerRepository) GetByID(id uuid.UUID) (*model.Trigger, error) {
	var trigger model.Trigger
	err := r.db.First(&trigger, "id = ?", id).Error
	return &trigger, err
}

// Create creates a new trigger
func (r *TriggerRepository) Create(trigger *model.Trigger) error {
	return r.db.Create(trigger).Error
}

// Update updates a trigger
func (r *TriggerRepository) Update(trigger *model.Trigger) error {
	return r.db.Save(trigger).Error
}

// Delete deletes a trigger
func (r *TriggerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Trigger{}, "id = ?", id).Error
}

// ArtifactRepository handles artifact database operations
type ArtifactRepository struct {
	db *gorm.DB
}

// NewArtifactRepository creates a new artifact repository
func NewArtifactRepository(db *gorm.DB) *ArtifactRepository {
	return &ArtifactRepository{db}
}

// ListByBuildID retrieves all artifacts for a build
func (r *ArtifactRepository) ListByBuildID(buildID uuid.UUID) ([]model.Artifact, error) {
	var artifacts []model.Artifact
	err := r.db.Where("build_id = ?", buildID).Find(&artifacts).Error
	return artifacts, err
}

// Create creates a new artifact
func (r *ArtifactRepository) Create(artifact *model.Artifact) error {
	return r.db.Create(artifact).Error
}

// CredentialRepository handles credential database operations
type CredentialRepository struct {
	db *gorm.DB
}

// NewCredentialRepository creates a new credential repository
func NewCredentialRepository(db *gorm.DB) *CredentialRepository {
	return &CredentialRepository{db}
}

// ListByUserID retrieves all credentials for a user
func (r *CredentialRepository) ListByUserID(userID uuid.UUID) ([]model.Credential, error) {
	var credentials []model.Credential
	err := r.db.Where("user_id = ?", userID).Find(&credentials).Error
	return credentials, err
}

// Create creates a new credential
func (r *CredentialRepository) Create(credential *model.Credential) error {
	return r.db.Create(credential).Error
}

// EnvVariableRepository handles environment variable database operations
type EnvVariableRepository struct {
	db *gorm.DB
}

// NewEnvVariableRepository creates a new env variable repository
func NewEnvVariableRepository(db *gorm.DB) *EnvVariableRepository {
	return &EnvVariableRepository{db}
}

// ListByPipelineID retrieves all env variables for a pipeline
func (r *EnvVariableRepository) ListByPipelineID(pipelineID uuid.UUID) ([]model.EnvVariable, error) {
	var vars []model.EnvVariable
	err := r.db.Where("pipeline_id = ?", pipelineID).Find(&vars).Error
	return vars, err
}

// Create creates a new env variable
func (r *EnvVariableRepository) Create(envVar *model.EnvVariable) error {
	return r.db.Create(envVar).Error
}
