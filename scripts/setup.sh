#!/bin/bash

# Development Setup Script for Portfolio Backend
# This script sets up the development environment

set -e

echo "🚀 Setting up Portfolio Backend development environment..."

# Check if required tools are installed
check_tool() {
    if ! command -v $1 &> /dev/null; then
        echo "❌ $1 is not installed. Please install it first."
        exit 1
    fi
}

echo "🔍 Checking required tools..."
check_tool "go"
check_tool "docker"
check_tool "docker-compose"
check_tool "git"

# Check Go version
GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | cut -c3-)
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" = "$REQUIRED_VERSION" ]; then
    echo "✅ Go version $GO_VERSION is compatible"
else
    echo "❌ Go version $GO_VERSION is not compatible. Required: $REQUIRED_VERSION or higher"
    exit 1
fi

# Create local environment file if it doesn't exist
if [ ! -f ".env" ]; then
    echo "📋 Creating local environment file..."
    cp .env.example .env
    echo "⚠️  Please edit .env file with your configuration"
fi

# Download Go dependencies
echo "📦 Downloading Go dependencies..."
go mod download
go mod tidy

# Install development tools
echo "🛠️  Installing development tools..."
go install github.com/air-verse/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest

# Start MongoDB with Docker Compose
echo "🐳 Starting MongoDB container..."
docker-compose up -d mongo

# Wait for MongoDB to be ready
echo "⏳ Waiting for MongoDB to be ready..."
sleep 10

# Run database migrations/initialization
echo "🗃️  Initializing database..."
mongosh mongodb://admin:password@localhost:27017/portfolio?authSource=admin < scripts/mongo-init.js

# Run tests
echo "🧪 Running tests..."
go test -v ./...

# Build the application
echo "🔨 Building application..."
go build -o portfolio-backend .

echo "✅ Development environment setup complete!"
echo ""
echo "🚀 To start the development server:"
echo "   make dev"
echo ""
echo "📚 Other useful commands:"
echo "   make test       - Run tests"
echo "   make lint       - Run linter"
echo "   make build      - Build application"
echo "   make docker     - Build Docker image"
echo "   make clean      - Clean build artifacts"
echo ""
echo "📖 Documentation:"
echo "   http://localhost:8080/docs    - API documentation"
echo "   http://localhost:8080/health  - Health check"