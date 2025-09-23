# Build stage
FROM golang:1.21-alpine AS builder

# Set build arguments
ARG VERSION=1.0.0
ARG BUILD_TIME
ARG GIT_COMMIT

# Install git for dependency management
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}" \
    -o portfolio-backend .

# Production stage
FROM alpine:3.18

# Install ca-certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 appuser && \
    adduser -u 1001 -G appuser -s /bin/sh -D appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/portfolio-backend .

# Copy environment example file
COPY --from=builder /app/.env.example .

# Change ownership to appuser
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8080

# Run the application
CMD ["./portfolio-backend"]