# Quick Start Guide

This guide will help you get the CI/CD Pipeline Builder up and running locally.

## Prerequisites

- Docker & Docker Compose
- Git
- A terminal/command line

## Local Development with Docker Compose

### 1. Clone and Setup

```bash
cd "Cicd pipeline builder"
cp .env.example .env
```

### 2. Build and Start Services

```bash
docker-compose up --build -d
```

This will start:
- PostgreSQL database (port 5432)
- Redis cache (port 6379)
- Go API server (port 8080)
- React frontend (port 3000)
- Job processor (background worker)

### 3. Initialize Database

Wait for database to be ready, then run migrations:

```bash
# Check database status
docker-compose logs postgres

# Run schema
docker exec cicd-db psql -U cicd -d cicd_db -f /docker-entrypoint-initdb.d/01-schema.sql
```

### 4. Access the Application

- **Frontend**: http://localhost:3000
- **API**: http://localhost:8080
- **API Health**: http://localhost:8080/health

### 5. Login

Default credentials (change in production):
- Username: `admin`
- Password: `admin123`

## Development Setup (Without Docker)

### Backend Development

```bash
cd backend

# Setup environment
cp ../.env.example .env

# Install dependencies
go mod download

# Run server
go run main.go
```

The API will start on port 8080.

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev
```

The frontend will start on http://localhost:5173 (Vite dev server)

### Job Processor Development

```bash
cd processor

# Install dependencies
go mod download

# Run processor
go run main.go
```

## Testing

### Backend Tests

```bash
cd backend
go test ./...
go test -v ./...  # Verbose
go test -cover ./...  # With coverage
```

### Frontend Tests

```bash
cd frontend
npm test
```

## Creating Your First Pipeline

1. Go to http://localhost:3000
2. Login with default credentials
3. Click "New Pipeline"
4. Fill in:
   - **Name**: My First Pipeline
   - **Repository URL**: https://github.com/user/repo
   - **Config Path**: .cicd.yml (default)
5. Click Create

## Setting Up GitHub Integration

### Create GitHub Token

1. Go to https://github.com/settings/tokens
2. Create new token with `repo` and `admin:repo_hook` scopes
3. Copy the token

### Configure Webhook

1. Go to your GitHub repository settings
2. Navigate to "Webhooks"
3. Add new webhook:
   - **Payload URL**: `http://your-domain:8080/api/v1/public/webhooks/github`
   - **Content type**: `application/json`
   - **Events**: Select "Push events" and "Pull request events"
   - **Secret**: Use a secure random value

4. Save webhook

### Test Webhook

Push to your repository and watch the build trigger automatically!

## Stopping Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## Troubleshooting

### Database Connection Error

```bash
# Check if database is running
docker-compose logs postgres

# Restart database
docker-compose restart postgres
```

### API Not Responding

```bash
# Check API logs
docker-compose logs api

# Restart API
docker-compose restart api
```

### Redis Connection Error

```bash
# Check Redis
docker-compose logs redis

# Restart Redis
docker-compose restart redis
```

### Port Already in Use

If port 8080, 3000, 5432, or 6379 is already in use:

1. Find process using the port:
   ```bash
   # Windows
   netstat -ano | findstr :8080
   
   # macOS/Linux
   lsof -i :8080
   ```

2. Kill the process or change Docker Compose ports in `docker-compose.yml`

### Logs Not Showing

```bash
# View all container logs
docker-compose logs

# Follow logs
docker-compose logs -f

# Specific service
docker-compose logs api
docker-compose logs processor
```

## Performance Tuning

### For Development

Default settings in `docker-compose.yml` are fine for development.

### For Production

See [KUBERNETES_DEPLOYMENT.md](KUBERNETES_DEPLOYMENT.md) for production setup recommendations.

## Next Steps

- [ ] Setup GitHub/GitLab integration
- [ ] Configure email notifications
- [ ] Setup Slack integration
- [ ] Configure Docker registry credentials
- [ ] Create scheduled triggers
- [ ] Deploy to production

## Additional Resources

- [API Documentation](API_DOCUMENTATION.md)
- [Kubernetes Deployment](KUBERNETES_DEPLOYMENT.md)
- [Example Pipeline Config](.cicd.yml.example)

## Support

For issues:
1. Check the logs: `docker-compose logs`
2. Check the GitHub issues
3. Create a new issue with:
   - Docker version
   - Docker Compose version
   - Operating system
   - Steps to reproduce
   - Error messages
