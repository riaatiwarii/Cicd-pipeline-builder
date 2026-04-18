# CI/CD Pipeline Builder

A scalable, production-ready CI/CD platform similar to Jenkins/GitLab CI. Built with Go, PostgreSQL, Redis, and React.

## 🏗️ Architecture Overview

```
┌─────────────┐        ┌──────────────┐        ┌─────────────┐
│   React     │◄──────►│  Go API      │◄──────►│ PostgreSQL  │
│  Frontend   │        │  Server      │        │             │
└─────────────┘        └──────────────┘        └─────────────┘
                             ▼
                       ┌──────────────┐
                       │    Redis     │
                       │  (Queue)     │
                       └──────────────┘
                             ▼
                       ┌──────────────┐
                       │ Job          │
                       │ Processors   │
                       │ (Workers)    │
                       └──────────────┘
```

## 📊 Tech Stack

### Backend
- **Framework**: Go + Gin/Echo
- **Database**: PostgreSQL with GORM
- **Queue**: Redis Streams
- **Container**: Docker SDK
- **Auth**: JWT

### Frontend
- **Framework**: React + TypeScript
- **State**: Redux Toolkit
- **UI**: Material-UI / Tailwind CSS
- **Real-time**: WebSocket + Socket.io
- **Build**: Vite

### Infrastructure
- **Containerization**: Docker
- **Orchestration**: Docker Compose (local), Kubernetes (production)
- **Message Queue**: Redis
- **Database**: PostgreSQL

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for frontend development)
- PostgreSQL 15+

### Setup

1. **Clone and setup**
```bash
cd "Cicd pipeline builder"
cp .env.example .env
```

2. **Start with Docker Compose**
```bash
docker-compose up -d
```

3. **Access the application**
- Frontend: http://localhost:3000
- API: http://localhost:8080
- API Docs: http://localhost:8080/swagger

4. **Create your first pipeline**
```bash
curl -X POST http://localhost:8080/api/v1/pipelines \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "My Pipeline",
    "repo_url": "https://github.com/user/repo",
    "config_path": ".cicd.yml"
  }'
```

## 📁 Project Structure

```
.
├── backend/                  # Go API Server
│   ├── main.go
│   ├── handler/             # HTTP handlers
│   ├── model/               # Data models
│   ├── service/             # Business logic
│   ├── repository/          # Database access
│   ├── middleware/          # Auth, logging, etc
│   ├── Dockerfile
│   └── go.mod
├── processor/               # Job Processor (Worker)
│   ├── main.go
│   ├── executor/            # Job execution logic
│   ├── docker/              # Docker integration
│   ├── logger/              # Logging
│   ├── Dockerfile
│   └── go.mod
├── frontend/                # React Dashboard
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── pages/          # Page components
│   │   ├── store/          # Redux store
│   │   ├── api/            # API client
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── package.json
│   ├── Dockerfile
│   └── vite.config.ts
├── database/                # Database schemas
│   ├── schema.sql          # Database DDL
│   └── migrations/         # Migration scripts
├── docker/                  # Docker configs
│   └── Dockerfile.*        # Various Dockerfiles
└── docker-compose.yml       # Local development setup
```

## 🔧 API Endpoints

### Pipelines
- `GET /api/v1/pipelines` - List all pipelines
- `POST /api/v1/pipelines` - Create pipeline
- `GET /api/v1/pipelines/:id` - Get pipeline details
- `PUT /api/v1/pipelines/:id` - Update pipeline
- `DELETE /api/v1/pipelines/:id` - Delete pipeline

### Builds/Jobs
- `GET /api/v1/pipelines/:id/builds` - List builds
- `POST /api/v1/pipelines/:id/trigger` - Trigger build
- `GET /api/v1/builds/:id` - Get build details
- `GET /api/v1/builds/:id/logs` - Stream build logs

### Webhooks
- `POST /api/v1/webhooks/github` - GitHub webhook
- `POST /api/v1/webhooks/gitlab` - GitLab webhook

### Triggers
- `GET /api/v1/triggers` - List triggers
- `POST /api/v1/triggers` - Create trigger
- `PUT /api/v1/triggers/:id` - Update trigger

## 🔐 Authentication

JWT-based authentication. Get token:
```bash
POST /api/v1/auth/login
{
  "username": "admin",
  "password": "password"
}
```

## 🐳 Pipeline YAML Config

```yaml
name: Build and Test
stages:
  - name: test
    image: golang:1.21
    script: |
      go test ./...
      go vet ./...
    timeout: 5m
    on_failure: continue

  - name: build
    image: golang:1.21
    script: |
      go build -o app
      ./app --version
    timeout: 10m
    artifacts:
      - app

  - name: push-docker
    image: docker:latest
    script: |
      docker build -t myapp:latest .
      docker push registry.io/myapp:latest
    timeout: 15m
    needs_docker: true

notifications:
  slack: "#ci-builds"
  email: "admin@example.com"
```

## 📝 Key Features

- ✅ **Multi-trigger support**: Webhooks, schedules, manual triggers
- ✅ **Real-time logs**: Live streaming via WebSocket
- ✅ **Docker isolation**: Each job runs in isolated container
- ✅ **Horizontal scaling**: Multiple workers for parallel execution
- ✅ **Artifact storage**: Save and download build artifacts
- ✅ **GitHub/GitLab integration**: Auto-trigger on push
- ✅ **Scheduled builds**: Cron-based pipeline scheduling
- ✅ **Build artifacts**: Store and retrieve build outputs
- ✅ **Status notifications**: Slack, email, webhooks
- ✅ **Pipeline visualization**: Real-time build status and logs

## 🏃 Development

### Backend Development
```bash
cd backend
go mod download
go run main.go
```

### Frontend Development
```bash
cd frontend
npm install
npm run dev
```

### Processor Development
```bash
cd processor
go mod download
go run main.go
```

## 🧪 Testing

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

## 📚 Database Schema

See [database/schema.sql](database/schema.sql) for complete schema.

Key tables:
- `pipelines` - Pipeline definitions
- `stages` - Pipeline stages/steps
- `builds` - Build execution records
- `jobs` - Individual job runs
- `logs` - Build logs
- `webhooks` - Webhook configurations
- `triggers` - Trigger definitions
- `artifacts` - Build artifacts

## 🚢 Deployment

### Docker Compose (Development)
```bash
docker-compose up -d
```

### Docker (Production)
```bash
docker build -t cicd-api ./backend
docker build -t cicd-processor ./processor
docker build -t cicd-frontend ./frontend

# Run with proper env vars and volume mounts
```

### Kubernetes
See `k8s/` directory for Kubernetes manifests.

## 📊 Monitoring

- Prometheus metrics on `/metrics`
- Health check on `/health`
- Logs in `logs/` directory
- Database logs in PostgreSQL

## 🤝 Contributing

1. Create a feature branch
2. Make your changes
3. Add tests
4. Submit pull request

## 📄 License

MIT

## 🆘 Troubleshooting

### Can't connect to database
- Ensure PostgreSQL is running
- Check DATABASE_URL in .env
- Run migrations: `go run database/migrate.go`

### API not responding
- Check if port 8080 is available
- View logs: `docker logs cicd-api`
- Ensure Redis is running

### Jobs not processing
- Check processor container is running
- Verify Redis connection
- Check job logs in database

## 📞 Support

For issues and questions, please open an issue on GitHub.
