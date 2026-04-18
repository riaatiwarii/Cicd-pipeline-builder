# CI/CD Pipeline Builder - API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All API requests (except login/register) require JWT authentication. Include the token in the `Authorization` header:

```
Authorization: Bearer <token>
```

## Endpoints

### Authentication

#### Login
```
POST /auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Register
```
POST /auth/register
Content-Type: application/json

{
  "username": "newuser",
  "email": "user@example.com",
  "password": "securepassword",
  "full_name": "Full Name"
}

Response:
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "newuser",
  "email": "user@example.com",
  "full_name": "Full Name",
  "role": "user",
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Pipelines

#### List Pipelines
```
GET /pipelines

Response:
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Go App Build",
    "description": "Build and test Go application",
    "repo_url": "https://github.com/user/repo",
    "repo_branch": "main",
    "config_path": ".cicd.yml",
    "status": "active",
    "is_public": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

#### Get Pipeline
```
GET /pipelines/{id}

Response: (same as above)
```

#### Create Pipeline
```
POST /pipelines
Content-Type: application/json

{
  "name": "Go App Build",
  "description": "Build and test Go application",
  "repo_url": "https://github.com/user/repo",
  "repo_branch": "main",
  "config_path": ".cicd.yml"
}

Response:
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  ...
}
```

#### Update Pipeline
```
PUT /pipelines/{id}
Content-Type: application/json

{
  "name": "Updated Pipeline Name",
  "status": "inactive"
}

Response: (same as get)
```

#### Delete Pipeline
```
DELETE /pipelines/{id}

Response: 204 No Content
```

#### Trigger Pipeline
```
POST /pipelines/{id}/trigger

Response:
{
  "id": "650e8400-e29b-41d4-a716-446655440000",
  "pipeline_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "pending",
  "trigger_type": "manual",
  "triggered_by": "550e8400-e29b-41d4-a716-446655440001",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Builds

#### List Builds for Pipeline
```
GET /pipelines/{id}/builds?limit=10&offset=0

Response:
[
  {
    "id": "650e8400-e29b-41d4-a716-446655440000",
    "pipeline_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "success",
    "branch": "main",
    "commit_hash": "abc123def456",
    "commit_message": "Fix: issue",
    "trigger_type": "webhook",
    "triggered_by": "550e8400-e29b-41d4-a716-446655440001",
    "started_at": "2024-01-01T00:00:00Z",
    "completed_at": "2024-01-01T00:05:00Z",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

#### Get Build Details
```
GET /builds/{id}

Response:
{
  "id": "650e8400-e29b-41d4-a716-446655440000",
  "pipeline_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "success",
  "branch": "main",
  "commit_hash": "abc123def456",
  "jobs": [
    {
      "id": "750e8400-e29b-41d4-a716-446655440000",
      "build_id": "650e8400-e29b-41d4-a716-446655440000",
      "stage_id": "550e8400-e29b-41d4-a716-446655440002",
      "status": "success",
      "exit_code": 0,
      "started_at": "2024-01-01T00:00:00Z",
      "completed_at": "2024-01-01T00:02:00Z"
    }
  ]
}
```

#### Get Build Logs
```
GET /builds/{id}/logs

Response:
[
  {
    "id": 1,
    "job_id": "750e8400-e29b-41d4-a716-446655440000",
    "line_number": 1,
    "log_level": "info",
    "message": "Starting stage...",
    "timestamp": "2024-01-01T00:00:00Z"
  },
  {
    "id": 2,
    "job_id": "750e8400-e29b-41d4-a716-446655440000",
    "line_number": 2,
    "log_level": "info",
    "message": "Running tests...",
    "timestamp": "2024-01-01T00:00:01Z"
  }
]
```

#### Get Build Artifacts
```
GET /builds/{id}/artifacts

Response:
[
  {
    "id": "850e8400-e29b-41d4-a716-446655440000",
    "build_id": "650e8400-e29b-41d4-a716-446655440000",
    "job_id": "750e8400-e29b-41d4-a716-446655440000",
    "name": "app",
    "file_path": "/artifacts/app-abc123",
    "file_size": 45678,
    "mime_type": "application/octet-stream",
    "created_at": "2024-01-01T00:02:00Z"
  }
]
```

#### Cancel Build
```
POST /builds/{id}/cancel

Response:
{
  "status": "cancelled",
  ...
}
```

### Webhooks

#### List Webhooks
```
GET /webhooks

Response:
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440010",
    "pipeline_id": "550e8400-e29b-41d4-a716-446655440000",
    "provider": "github",
    "events": ["push", "pull_request"],
    "webhook_url": "https://yourdomain.com/webhooks/github",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

#### Create Webhook
```
POST /webhooks
Content-Type: application/json

{
  "pipeline_id": "550e8400-e29b-41d4-a716-446655440000",
  "provider": "github",
  "events": ["push", "pull_request"]
}

Response:
{
  "id": "550e8400-e29b-41d4-a716-446655440010",
  "pipeline_id": "550e8400-e29b-41d4-a716-446655440000",
  "provider": "github",
  "events": ["push", "pull_request"],
  "secret_key": "whsec_1234567890",
  "webhook_url": "https://yourdomain.com/webhooks/github",
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### GitHub Webhook (Public)
```
POST /public/webhooks/github
Content-Type: application/json
X-Hub-Signature-256: sha256=<signature>

{
  "repository": {
    "name": "repo-name",
    "url": "https://github.com/user/repo"
  },
  "branch": "main",
  "commits": [
    {
      "id": "abc123def456",
      "message": "Fix: issue"
    }
  ]
}

Response:
{
  "status": "received",
  "build_id": "650e8400-e29b-41d4-a716-446655440000"
}
```

### Triggers

#### List Triggers
```
GET /triggers

Response:
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440020",
    "pipeline_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Daily Build",
    "trigger_type": "schedule",
    "cron_expression": "0 0 * * *",
    "is_active": true,
    "last_triggered_at": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

#### Create Trigger
```
POST /triggers
Content-Type: application/json

{
  "pipeline_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Daily Build",
  "trigger_type": "schedule",
  "cron_expression": "0 0 * * *"
}

Response:
{
  "id": "550e8400-e29b-41d4-a716-446655440020",
  ...
}
```

## Error Responses

All errors follow this format:

```json
{
  "error": "Error message here",
  "code": "ERROR_CODE"
}
```

Common HTTP Status Codes:
- `200 OK`: Successful request
- `201 Created`: Resource created
- `202 Accepted`: Request accepted (async)
- `204 No Content`: Success, no response body
- `400 Bad Request`: Invalid request
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict
- `500 Internal Server Error`: Server error

## Rate Limiting

API rate limits:
- 1000 requests per hour for authenticated users
- 100 requests per hour for public endpoints

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200
```

## WebSocket Events (Real-time Logs)

Connect to WebSocket at `ws://localhost:8080/ws/builds/{buildId}/logs`

Message format:
```json
{
  "type": "log",
  "job_id": "750e8400-e29b-41d4-a716-446655440000",
  "message": "Test output...",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Examples

### Create and Trigger Pipeline

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}' | jq -r '.token')

# Create pipeline
curl -s -X POST http://localhost:8080/api/v1/pipelines \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Pipeline",
    "repo_url": "https://github.com/user/repo",
    "repo_branch": "main",
    "config_path": ".cicd.yml"
  }'

# Trigger pipeline
PIPELINE_ID="550e8400-e29b-41d4-a716-446655440000"
curl -s -X POST http://localhost:8080/api/v1/pipelines/$PIPELINE_ID/trigger \
  -H "Authorization: Bearer $TOKEN"
```

### Stream Build Logs

```bash
# Get build logs via HTTP
curl -s http://localhost:8080/api/v1/builds/650e8400-e29b-41d4-a716-446655440000/logs \
  -H "Authorization: Bearer $TOKEN" | jq '.[]'

# Or connect to WebSocket for real-time logs
wscat -c ws://localhost:8080/ws/builds/650e8400-e29b-41d4-a716-446655440000/logs
```
