# Build stage
FROM golang:1.25.1-alpine AS builder

ARG GO_VERSION=1.25.1
ARG BUILD_ENV=production

WORKDIR /app

# Copy go mod files from api directory (build context is repo root)
COPY api/go.mod api/go.sum ./
RUN go mod download

# Copy api source code
COPY api/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose API port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]

