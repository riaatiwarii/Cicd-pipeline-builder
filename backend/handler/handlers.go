package handler

import (
	"net/http"
	"strconv"

	"cicd-pipeline-builder/backend/model"
	"cicd-pipeline-builder/backend/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(req.Username, req.Email, req.Password, req.FullName)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// RefreshToken refreshes JWT token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Implementation would verify old token and issue new one
	c.JSON(http.StatusOK, gin.H{"message": "token refreshed"})
}

// PipelineHandler handles pipeline requests
type PipelineHandler struct {
	pipelineService *service.PipelineService
	buildService    *service.BuildService
}

// NewPipelineHandler creates a new pipeline handler
func NewPipelineHandler(pipelineService *service.PipelineService, buildService *service.BuildService) *PipelineHandler {
	return &PipelineHandler{pipelineService, buildService}
}

// ListPipelines lists all pipelines for the user
func (h *PipelineHandler) ListPipelines(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	pipelines, err := h.pipelineService.ListPipelines(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pipelines)
}

// GetPipeline retrieves a specific pipeline
func (h *PipelineHandler) GetPipeline(c *gin.Context) {
	id := c.Param("id")
	pipelineID, _ := uuid.Parse(id)

	pipeline, err := h.pipelineService.GetPipeline(pipelineID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pipeline not found"})
		return
	}

	c.JSON(http.StatusOK, pipeline)
}

// CreatePipeline creates a new pipeline
func (h *PipelineHandler) CreatePipeline(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var pipeline model.Pipeline
	if err := c.BindJSON(&pipeline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pipeline.UserID = uid
	if err := h.pipelineService.CreatePipeline(&pipeline); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pipeline)
}

// UpdatePipeline updates a pipeline
func (h *PipelineHandler) UpdatePipeline(c *gin.Context) {
	id := c.Param("id")
	pipelineID, _ := uuid.Parse(id)

	var pipeline model.Pipeline
	if err := c.BindJSON(&pipeline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pipeline.ID = pipelineID
	if err := h.pipelineService.UpdatePipeline(&pipeline); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pipeline)
}

// DeletePipeline deletes a pipeline
func (h *PipelineHandler) DeletePipeline(c *gin.Context) {
	id := c.Param("id")
	pipelineID, _ := uuid.Parse(id)

	if err := h.pipelineService.DeletePipeline(pipelineID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// TriggerPipeline manually triggers a pipeline
func (h *PipelineHandler) TriggerPipeline(c *gin.Context) {
	id := c.Param("id")
	pipelineID, _ := uuid.Parse(id)
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	build, err := h.pipelineService.TriggerPipeline(pipelineID, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, build)
}

// BuildHandler handles build requests
type BuildHandler struct {
	buildService *service.BuildService
}

// NewBuildHandler creates a new build handler
func NewBuildHandler(buildService *service.BuildService) *BuildHandler {
	return &BuildHandler{buildService}
}

// ListBuilds lists builds for a pipeline
func (h *BuildHandler) ListBuilds(c *gin.Context) {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := c.Query("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}

	pipelineID := c.Param("id")
	pid, _ := uuid.Parse(pipelineID)

	builds, err := h.buildService.ListBuilds(pid, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, builds)
}

// GetPipelineBuilds lists builds for a pipeline
func (h *BuildHandler) GetPipelineBuilds(c *gin.Context) {
	h.ListBuilds(c)
}

// GetBuild retrieves a specific build
func (h *BuildHandler) GetBuild(c *gin.Context) {
	id := c.Param("id")
	buildID, _ := uuid.Parse(id)

	build, err := h.buildService.GetBuild(buildID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "build not found"})
		return
	}

	c.JSON(http.StatusOK, build)
}

// GetBuildLogs retrieves logs for a build
func (h *BuildHandler) GetBuildLogs(c *gin.Context) {
	id := c.Param("id")
	buildID, _ := uuid.Parse(id)

	logs, err := h.buildService.GetBuildLogs(buildID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetBuildArtifacts retrieves artifacts for a build
func (h *BuildHandler) GetBuildArtifacts(c *gin.Context) {
	c.JSON(http.StatusOK, []interface{}{})
}

// CancelBuild cancels a build
func (h *BuildHandler) CancelBuild(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "build cancelled"})
}

// WebhookHandler handles webhook requests
type WebhookHandler struct {
	webhookService *service.WebhookService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(webhookService *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{webhookService}
}

// ListWebhooks lists webhooks
func (h *WebhookHandler) ListWebhooks(c *gin.Context) {
	c.JSON(http.StatusOK, []interface{}{})
}

// CreateWebhook creates a new webhook
func (h *WebhookHandler) CreateWebhook(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"id": "webhook-123"})
}

// GetWebhook retrieves a webhook
func (h *WebhookHandler) GetWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// UpdateWebhook updates a webhook
func (h *WebhookHandler) UpdateWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteWebhook deletes a webhook
func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

// HandleGitHubWebhook handles GitHub webhooks
func (h *WebhookHandler) HandleGitHubWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "webhook received"})
}

// HandleGitLabWebhook handles GitLab webhooks
func (h *WebhookHandler) HandleGitLabWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "webhook received"})
}

// TriggerHandler handles trigger requests
type TriggerHandler struct {
	triggerService *service.TriggerService
}

// NewTriggerHandler creates a new trigger handler
func NewTriggerHandler(triggerService *service.TriggerService) *TriggerHandler {
	return &TriggerHandler{triggerService}
}

// ListTriggers lists triggers
func (h *TriggerHandler) ListTriggers(c *gin.Context) {
	c.JSON(http.StatusOK, []interface{}{})
}

// CreateTrigger creates a new trigger
func (h *TriggerHandler) CreateTrigger(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"id": "trigger-123"})
}

// GetTrigger retrieves a trigger
func (h *TriggerHandler) GetTrigger(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// UpdateTrigger updates a trigger
func (h *TriggerHandler) UpdateTrigger(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteTrigger deletes a trigger
func (h *TriggerHandler) DeleteTrigger(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

// HealthHandler handles health checks
type HealthHandler struct {
	db *gorm.DB
	redis interface{}
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB, redis interface{}) *HealthHandler {
	return &HealthHandler{db, redis}
}

// Health returns the health status
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

// Ready returns the ready status
func (h *HealthHandler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
