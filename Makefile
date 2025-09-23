# Portfolio Backend Makefile
# Provides convenient commands for development and deployment

.PHONY: help dev build test lint clean docker run deps setup

# Default target
help: ## Show this help message
	@echo "Portfolio Backend - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development
dev: ## Start development server with hot reload
	@echo "🚀 Starting development server..."
	air

setup: ## Setup development environment
	@echo "🛠️  Setting up development environment..."
	./scripts/setup.sh

deps: ## Download and tidy dependencies
	@echo "📦 Managing dependencies..."
	go mod download
	go mod tidy

# Building
build: ## Build the application
	@echo "🔨 Building application..."
	CGO_ENABLED=0 GOOS=linux go build \
		-ldflags="-w -s -X main.version=$$(git describe --tags --always) -X main.buildTime=$$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.gitCommit=$$(git rev-parse HEAD)" \
		-o portfolio-backend .

build-windows: ## Build for Windows
	@echo "🔨 Building for Windows..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
		-ldflags="-w -s -X main.version=$$(git describe --tags --always) -X main.buildTime=$$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.gitCommit=$$(git rev-parse HEAD)" \
		-o portfolio-backend.exe .

build-mac: ## Build for macOS
	@echo "🔨 Building for macOS..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
		-ldflags="-w -s -X main.version=$$(git describe --tags --always) -X main.buildTime=$$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.gitCommit=$$(git rev-parse HEAD)" \
		-o portfolio-backend-mac .

# Testing
test: ## Run all tests
	@echo "🧪 Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: ## Run tests with coverage report
	@echo "🧪 Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

benchmark: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	go test -bench=. -benchmem ./...

# Code Quality
lint: ## Run linter
	@echo "🔍 Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "🎨 Formatting code..."
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	go vet ./...

# Docker
docker: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build \
		--build-arg VERSION=$$(git describe --tags --always) \
		--build-arg BUILD_TIME=$$(date -u +%Y-%m-%dT%H:%M:%SZ) \
		--build-arg GIT_COMMIT=$$(git rev-parse HEAD) \
		-t portfolio-backend:latest .

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run --rm -p 8080:8080 --env-file .env portfolio-backend:latest

docker-compose-up: ## Start all services with docker-compose
	@echo "🐳 Starting all services..."
	docker-compose up -d

docker-compose-down: ## Stop all services
	@echo "🐳 Stopping all services..."
	docker-compose down

docker-compose-logs: ## View docker-compose logs
	@echo "📋 Viewing logs..."
	docker-compose logs -f

# Database
db-start: ## Start MongoDB container
	@echo "🗃️  Starting MongoDB..."
	docker-compose up -d mongo

db-stop: ## Stop MongoDB container
	@echo "🗃️  Stopping MongoDB..."
	docker-compose stop mongo

db-init: ## Initialize database with sample data
	@echo "🗃️  Initializing database..."
	mongosh mongodb://admin:password@localhost:27017/portfolio?authSource=admin < scripts/mongo-init.js

db-backup: ## Backup database
	@echo "💾 Creating database backup..."
	docker exec -t $$(docker-compose ps -q mongo) mongodump --host localhost --port 27017 --username admin --password password --authenticationDatabase admin --db portfolio --out /tmp/backup
	docker cp $$(docker-compose ps -q mongo):/tmp/backup ./backups/$$(date +%Y%m%d_%H%M%S)

# Running
run: build ## Build and run the application
	@echo "🚀 Starting application..."
	./portfolio-backend

run-prod: ## Run in production mode
	@echo "🚀 Starting in production mode..."
	GIN_MODE=release ./portfolio-backend

# Deployment
deploy-local: ## Deploy locally using docker-compose
	@echo "🚀 Deploying locally..."
	./deploy.sh local

deploy-staging: ## Deploy to staging environment
	@echo "🚀 Deploying to staging..."
	./deploy.sh staging

deploy-prod: ## Deploy to production environment
	@echo "🚀 Deploying to production..."
	./deploy.sh production

# Utilities
clean: ## Clean build artifacts and dependencies
	@echo "🧹 Cleaning..."
	go clean
	rm -f portfolio-backend portfolio-backend.exe portfolio-backend-mac
	rm -f coverage.out coverage.html
	docker system prune -f

docs: ## Generate API documentation
	@echo "📚 Generating documentation..."
	swag init

health: ## Check application health
	@echo "🏥 Checking health..."
	curl -f http://localhost:8080/health || echo "❌ Health check failed"

logs: ## View application logs
	@echo "📋 Viewing logs..."
	docker-compose logs -f portfolio-backend

# Git hooks
install-hooks: ## Install git hooks
	@echo "🪝 Installing git hooks..."
	cp scripts/pre-commit .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

# Security
security-scan: ## Run security scan with gosec
	@echo "🔒 Running security scan..."
	gosec ./...

audit: ## Run dependency audit
	@echo "🔍 Running dependency audit..."
	go list -json -m all | nancy sleuth

# Release
tag: ## Create a new git tag
	@read -p "Enter tag version (e.g., v1.0.0): " VERSION; \
	git tag -a $$VERSION -m "Release $$VERSION"; \
	echo "✅ Tag $$VERSION created. Push with: git push origin $$VERSION"

release: test lint build docker ## Run full release process
	@echo "🎉 Release process completed!"