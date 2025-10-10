# Build stage
FROM golang:1.25.1-alpine AS builder

ARG GO_VERSION=1.25.1
ARG BUILD_ENV=production

WORKDIR /app

# Install build dependencies for CGO (required for go-sqlite3)
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod files from api directory (build context is repo root)
COPY api/go.mod api/go.sum ./
RUN go mod download

# Copy api source code
COPY api/ ./

# Build the application with CGO enabled (required for go-sqlite3)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /app/server ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies for SQLite
RUN apk add --no-cache sqlite-libs

# Copy binary from builder
COPY --from=builder /app/server .

# Expose API port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/healthz || exit 1

CMD ["./server"]

