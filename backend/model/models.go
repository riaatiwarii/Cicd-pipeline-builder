package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// User represents a system user
type User struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:255" json:"username"`
	Email     string    `gorm:"uniqueIndex;size:255" json:"email"`
	Password  string    `json:"-"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"` // admin, user
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin *time.Time `json:"last_login"`
}

// Pipeline represents a CI/CD pipeline
type Pipeline struct {
	ID         uuid.UUID      `gorm:"primaryKey" json:"id"`
	UserID     uuid.UUID      `json:"user_id"`
	User       *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name       string         `gorm:"uniqueIndex:idx_user_pipeline;size:255" json:"name"`
	Description string        `json:"description"`
	RepoURL    string         `json:"repo_url"`
	RepoBranch string         `json:"repo_branch" gorm:"default:main"`
	ConfigPath string         `json:"config_path" gorm:"default:.cicd.yml"`
	Status     string         `json:"status" gorm:"default:active"` // active, inactive, archived
	IsPublic   bool           `json:"is_public" gorm:"default:false"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	
	// Relations
	Stages      []Stage         `json:"stages,omitempty" gorm:"foreignKey:PipelineID"`
	Builds      []Build         `json:"builds,omitempty" gorm:"foreignKey:PipelineID"`
	Webhooks    []Webhook       `json:"webhooks,omitempty" gorm:"foreignKey:PipelineID"`
	Triggers    []Trigger       `json:"triggers,omitempty" gorm:"foreignKey:PipelineID"`
	EnvVars     []EnvVariable   `json:"env_variables,omitempty" gorm:"foreignKey:PipelineID"`
}

// Stage represents a stage/step in a pipeline
type Stage struct {
	ID           uuid.UUID `gorm:"primaryKey" json:"id"`
	PipelineID   uuid.UUID `gorm:"uniqueIndex:idx_pipeline_stage;index" json:"pipeline_id"`
	Pipeline     *Pipeline `json:"pipeline,omitempty" gorm:"foreignKey:PipelineID"`
	Name         string    `gorm:"uniqueIndex:idx_pipeline_stage;size:255" json:"name"`
	StageOrder   int       `json:"stage_order"`
	Description  string    `json:"description"`
	DockerImage  string    `json:"docker_image" gorm:"default:ubuntu:22.04"`
	Script       string    `gorm:"type:text" json:"script"`
	TimeoutSeconds int    `json:"timeout_seconds" gorm:"default:3600"`
	AllowFailure bool     `json:"allow_failure" gorm:"default:false"`
	NeedsDocker  bool     `json:"needs_docker" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Relations
	Jobs []Job `json:"jobs,omitempty" gorm:"foreignKey:StageID"`
}

// Build represents a pipeline execution
type Build struct {
	ID          uuid.UUID  `gorm:"primaryKey" json:"id"`
	PipelineID  uuid.UUID  `gorm:"index" json:"pipeline_id"`
	Pipeline    *Pipeline  `json:"pipeline,omitempty" gorm:"foreignKey:PipelineID"`
	Status      string     `json:"status" gorm:"default:pending"` // pending, running, success, failed, cancelled
	Branch      string     `json:"branch"`
	CommitHash  string     `json:"commit_hash"`
	CommitMessage string   `gorm:"type:text" json:"commit_message"`
	TriggeredBy *uuid.UUID `json:"triggered_by"`
	TriggerType string     `json:"trigger_type"` // webhook, schedule, manual
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	
	// Relations
	Jobs      []Job      `json:"jobs,omitempty" gorm:"foreignKey:BuildID"`
	Artifacts []Artifact `json:"artifacts,omitempty" gorm:"foreignKey:BuildID"`
}

// Job represents a single stage execution within a build
type Job struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	BuildID     uuid.UUID `gorm:"index" json:"build_id"`
	Build       *Build    `json:"build,omitempty" gorm:"foreignKey:BuildID"`
	StageID     uuid.UUID `gorm:"index" json:"stage_id"`
	Stage       *Stage    `json:"stage,omitempty" gorm:"foreignKey:StageID"`
	Status      string    `json:"status" gorm:"default:pending"` // pending, running, success, failed, cancelled
	ExitCode    *int      `json:"exit_code"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	ContainerID string    `json:"container_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relations
	Logs      []Log      `json:"logs,omitempty" gorm:"foreignKey:JobID"`
	Artifacts []Artifact `json:"artifacts,omitempty" gorm:"foreignKey:JobID"`
}

// Log represents a log entry from a job
type Log struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	JobID     uuid.UUID `gorm:"index" json:"job_id"`
	Job       *Job      `json:"job,omitempty" gorm:"foreignKey:JobID"`
	LineNumber int      `json:"line_number"`
	LogLevel  string    `json:"log_level" gorm:"default:info"`
	Message   string    `gorm:"type:text" json:"message"`
	Timestamp time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"timestamp"`
}

// Webhook represents a webhook configuration
type Webhook struct {
	ID        uuid.UUID      `gorm:"primaryKey" json:"id"`
	PipelineID uuid.UUID     `gorm:"index" json:"pipeline_id"`
	Pipeline  *Pipeline      `json:"pipeline,omitempty" gorm:"foreignKey:PipelineID"`
	Provider  string         `json:"provider"` // github, gitlab, gitea, custom
	Events    datatypes.JSONSlice `gorm:"type:jsonb" json:"events"`
	SecretKey string         `json:"secret_key"`
	WebhookURL string        `json:"webhook_url"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	
	// Relations
	Events2 []WebhookEvent `json:"webhook_events,omitempty" gorm:"foreignKey:WebhookID"`
}

// Trigger represents a trigger configuration (schedule, manual, etc)
type Trigger struct {
	ID               uuid.UUID `gorm:"primaryKey" json:"id"`
	PipelineID       uuid.UUID `gorm:"index" json:"pipeline_id"`
	Pipeline         *Pipeline `json:"pipeline,omitempty" gorm:"foreignKey:PipelineID"`
	Name             string    `json:"name"`
	TriggerType      string    `json:"trigger_type"` // webhook, schedule, manual
	CronExpression   string    `json:"cron_expression"`
	IsActive         bool      `json:"is_active" gorm:"default:true"`
	LastTriggeredAt  *time.Time `json:"last_triggered_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Artifact represents a build artifact
type Artifact struct {
	ID        uuid.UUID  `gorm:"primaryKey" json:"id"`
	BuildID   uuid.UUID  `gorm:"index" json:"build_id"`
	Build     *Build     `json:"build,omitempty" gorm:"foreignKey:BuildID"`
	JobID     *uuid.UUID `json:"job_id"`
	Job       *Job       `json:"job,omitempty" gorm:"foreignKey:JobID"`
	Name      string     `json:"name"`
	FilePath  string     `json:"file_path"`
	FileSize  int64      `json:"file_size"`
	MimeType  string     `json:"mime_type"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// Credential represents an encrypted credential
type Credential struct {
	ID              uuid.UUID `gorm:"primaryKey" json:"id"`
	UserID          uuid.UUID `gorm:"index" json:"user_id"`
	User            *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CredentialType  string    `json:"credential_type"` // github_token, docker_registry, ssh_key, etc
	Name            string    `gorm:"uniqueIndex:idx_user_cred;size:255" json:"name"`
	EncryptedValue  string    `json:"-"`
	IsActive        bool      `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// EnvVariable represents an environment variable for a pipeline
type EnvVariable struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	PipelineID uuid.UUID `gorm:"index" json:"pipeline_id"`
	Pipeline  *Pipeline `json:"pipeline,omitempty" gorm:"foreignKey:PipelineID"`
	Name      string    `gorm:"uniqueIndex:idx_pipeline_env;size:255" json:"name"`
	Value     string    `gorm:"type:text" json:"value"`
	IsSecret  bool      `json:"is_secret" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WebhookEvent represents a webhook event record for audit trail
type WebhookEvent struct {
	ID        uuid.UUID      `gorm:"primaryKey" json:"id"`
	WebhookID uuid.UUID      `gorm:"index" json:"webhook_id"`
	Webhook   *Webhook       `json:"webhook,omitempty" gorm:"foreignKey:WebhookID"`
	BuildID   *uuid.UUID     `json:"build_id"`
	Build     *Build         `json:"build,omitempty" gorm:"foreignKey:BuildID"`
	EventType string         `json:"event_type"`
	Payload   datatypes.JSON `gorm:"type:jsonb" json:"payload"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
}

// TableName specifies custom table names
func (User) TableName() string { return "users" }
func (Pipeline) TableName() string { return "pipelines" }
func (Stage) TableName() string { return "stages" }
func (Build) TableName() string { return "builds" }
func (Job) TableName() string { return "jobs" }
func (Log) TableName() string { return "logs" }
func (Webhook) TableName() string { return "webhooks" }
func (Trigger) TableName() string { return "triggers" }
func (Artifact) TableName() string { return "artifacts" }
func (Credential) TableName() string { return "credentials" }
func (EnvVariable) TableName() string { return "env_variables" }
func (WebhookEvent) TableName() string { return "webhook_events" }
