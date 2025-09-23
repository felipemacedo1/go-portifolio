#!/bin/bash

# Portfolio Backend Deployment Script
# Usage: ./deploy.sh [environment]

set -e

ENVIRONMENT=${1:-production}
PROJECT_NAME="portfolio-backend"
REGISTRY="your-registry.com"  # Replace with your Docker registry

echo "🚀 Starting deployment for environment: $ENVIRONMENT"

# Load environment-specific configuration
if [ -f ".env.$ENVIRONMENT" ]; then
    echo "📋 Loading environment configuration..."
    source ".env.$ENVIRONMENT"
else
    echo "⚠️  Warning: .env.$ENVIRONMENT not found, using defaults"
fi

# Build Docker image
echo "🐳 Building Docker image..."
docker build \
    --build-arg VERSION=$(git describe --tags --always) \
    --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
    --build-arg GIT_COMMIT=$(git rev-parse HEAD) \
    -t $PROJECT_NAME:latest \
    -t $PROJECT_NAME:$(git describe --tags --always) \
    .

# Run tests
echo "🧪 Running tests..."
docker run --rm \
    -v $(pwd):/app \
    -w /app \
    golang:1.21-alpine \
    sh -c "go mod download && go test -v ./..."

# Security scan (if trivy is available)
if command -v trivy &> /dev/null; then
    echo "🔒 Running security scan..."
    trivy image $PROJECT_NAME:latest
fi

case $ENVIRONMENT in
    "local")
        echo "🏠 Deploying locally with docker-compose..."
        docker-compose down
        docker-compose up -d
        echo "✅ Local deployment complete!"
        echo "📱 Application available at: http://localhost:8080"
        ;;
    
    "staging"|"production")
        echo "☁️  Deploying to $ENVIRONMENT..."
        
        # Tag and push to registry
        docker tag $PROJECT_NAME:latest $REGISTRY/$PROJECT_NAME:latest
        docker tag $PROJECT_NAME:latest $REGISTRY/$PROJECT_NAME:$(git describe --tags --always)
        
        echo "📤 Pushing to registry..."
        docker push $REGISTRY/$PROJECT_NAME:latest
        docker push $REGISTRY/$PROJECT_NAME:$(git describe --tags --always)
        
        # Deploy using your preferred method (kubectl, helm, etc.)
        if command -v kubectl &> /dev/null; then
            echo "⚙️  Applying Kubernetes manifests..."
            kubectl apply -f k8s/
        fi
        
        echo "✅ $ENVIRONMENT deployment complete!"
        ;;
    
    *)
        echo "❌ Unknown environment: $ENVIRONMENT"
        echo "Available environments: local, staging, production"
        exit 1
        ;;
esac

# Health check
echo "🏥 Performing health check..."
sleep 10
if [ "$ENVIRONMENT" = "local" ]; then
    HEALTH_URL="http://localhost:8080/health"
else
    HEALTH_URL="$APP_URL/health"
fi

if curl -f -s $HEALTH_URL > /dev/null; then
    echo "✅ Health check passed!"
    echo "🎉 Deployment successful!"
else
    echo "❌ Health check failed!"
    exit 1
fi