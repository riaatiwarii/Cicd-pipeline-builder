.PHONY: help build up down logs clean test lint fmt install-tools

# Variables
BACKEND_DIR := backend
PROCESSOR_DIR := processor
FRONTEND_DIR := frontend
COMPOSE_FILE := docker-compose.yml

help: ## Show this help message
	@echo 'Usage: make [target]'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# Docker Compose
up: ## Start all services with Docker Compose
	docker-compose up -d

down: ## Stop all services
	docker-compose down

logs: ## View logs from all services
	docker-compose logs -f

logs-api: ## View API logs
	docker-compose logs -f api

logs-processor: ## View processor logs
	docker-compose logs -f processor

logs-frontend: ## View frontend logs
	docker-compose logs -f frontend

clean: ## Clean up containers and volumes
	docker-compose down -v

rebuild: ## Rebuild and restart services
	docker-compose down -v
	docker-compose up --build -d

# Backend
backend-build: ## Build backend binary
	cd $(BACKEND_DIR) && CGO_ENABLED=1 GOOS=linux go build -o app main.go

backend-run: ## Run backend locally
	cd $(BACKEND_DIR) && go run main.go

backend-test: ## Run backend tests
	cd $(BACKEND_DIR) && go test ./...

backend-test-v: ## Run backend tests with verbose output
	cd $(BACKEND_DIR) && go test -v ./...

backend-coverage: ## Run backend tests with coverage
	cd $(BACKEND_DIR) && go test -cover ./...

backend-lint: ## Lint backend code
	cd $(BACKEND_DIR) && golangci-lint run

backend-fmt: ## Format backend code
	cd $(BACKEND_DIR) && go fmt ./...

backend-tidy: ## Tidy backend dependencies
	cd $(BACKEND_DIR) && go mod tidy

# Processor
processor-build: ## Build processor binary
	cd $(PROCESSOR_DIR) && CGO_ENABLED=1 GOOS=linux go build -o processor main.go

processor-run: ## Run processor locally
	cd $(PROCESSOR_DIR) && go run main.go

processor-test: ## Run processor tests
	cd $(PROCESSOR_DIR) && go test ./...

processor-lint: ## Lint processor code
	cd $(PROCESSOR_DIR) && golangci-lint run

processor-fmt: ## Format processor code
	cd $(PROCESSOR_DIR) && go fmt ./...

processor-tidy: ## Tidy processor dependencies
	cd $(PROCESSOR_DIR) && go mod tidy

# Frontend
frontend-install: ## Install frontend dependencies
	cd $(FRONTEND_DIR) && npm install

frontend-run: ## Run frontend dev server
	cd $(FRONTEND_DIR) && npm run dev

frontend-build: ## Build frontend for production
	cd $(FRONTEND_DIR) && npm run build

frontend-test: ## Run frontend tests
	cd $(FRONTEND_DIR) && npm test

frontend-lint: ## Lint frontend code
	cd $(FRONTEND_DIR) && npm run lint

frontend-clean: ## Clean frontend dependencies
	rm -rf $(FRONTEND_DIR)/node_modules

# Testing
test: backend-test frontend-test ## Run all tests

test-integration: ## Run integration tests
	docker-compose up -d
	sleep 5
	cd $(BACKEND_DIR) && go test -tags=integration ./...
	docker-compose down

# Code Quality
lint: backend-lint frontend-lint ## Lint all code

fmt: backend-fmt frontend-lint ## Format all code

# Installation
install-go-tools: ## Install Go development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cosmtrek/air@latest

install-tools: install-go-tools frontend-install ## Install all development tools

# Database
db-shell: ## Open database shell
	docker exec -it cicd-db psql -U cicd -d cicd_db

db-backup: ## Backup database
	docker exec cicd-db pg_dump -U cicd -d cicd_db > backup_$$(date +%Y%m%d_%H%M%S).sql

db-reset: ## Reset database (careful!)
	docker-compose restart postgres

# Docker
docker-build: ## Build all Docker images
	docker build -t cicd-api:latest ./$(BACKEND_DIR)
	docker build -t cicd-processor:latest ./$(PROCESSOR_DIR)
	docker build -t cicd-frontend:latest ./$(FRONTEND_DIR)

docker-push: ## Push images to registry (set REGISTRY env var)
	docker tag cicd-api:latest $$REGISTRY/cicd-api:latest
	docker tag cicd-processor:latest $$REGISTRY/cicd-processor:latest
	docker tag cicd-frontend:latest $$REGISTRY/cicd-frontend:latest
	docker push $$REGISTRY/cicd-api:latest
	docker push $$REGISTRY/cicd-processor:latest
	docker push $$REGISTRY/cicd-frontend:latest

# Development
dev-backend: ## Run backend in dev mode (with live reload)
	cd $(BACKEND_DIR) && air

dev-frontend: ## Run frontend in dev mode
	cd $(FRONTEND_DIR) && npm run dev

dev-processor: ## Run processor in dev mode
	cd $(PROCESSOR_DIR) && go run main.go

# Documentation
docs-serve: ## Serve documentation
	python -m http.server 8000

# Cleanup
clean-all: clean ## Full cleanup
	rm -rf $(BACKEND_DIR)/app
	rm -rf $(PROCESSOR_DIR)/processor
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules
	rm -f backup_*.sql

# Setup
setup: install-tools backend-tidy processor-tidy ## Initial setup
	@echo "Setup complete! Run 'make up' to start services"

# CI/CD Helpers
verify: lint test ## Verify all code quality checks pass

.DEFAULT_GOAL := help
