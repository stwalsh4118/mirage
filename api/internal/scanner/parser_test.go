package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDockerfile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "exact Dockerfile in root",
			path:     "Dockerfile",
			expected: true,
		},
		{
			name:     "Dockerfile in subdirectory",
			path:     "services/api/Dockerfile",
			expected: true,
		},
		{
			name:     "dockerfile with .dockerfile extension",
			path:     "services/api/prod.dockerfile",
			expected: true,
		},
		{
			name:     "Dockerfile with suffix not .dockerfile",
			path:     "apps/frontend/Dockerfile.prod",
			expected: false, // Only exact "Dockerfile" or "*.dockerfile" pattern
		},
		{
			name:     "not a dockerfile",
			path:     "package.json",
			expected: false,
		},
		{
			name:     "dockerfile in filename but not extension",
			path:     "dockerfile-template.txt",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDockerfile(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInferServiceName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple subdirectory",
			path:     "services/api/Dockerfile",
			expected: "api",
		},
		{
			name:     "nested subdirectory",
			path:     "packages/backend/services/auth/Dockerfile",
			expected: "auth",
		},
		{
			name:     "root directory",
			path:     "Dockerfile",
			expected: "root",
		},
		{
			name:     "with .dockerfile extension",
			path:     "services/worker/prod.dockerfile",
			expected: "worker",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inferServiceName(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseDockerfile(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectedPorts []int
		expectedArgs  []string
		expectedBase  string
	}{
		{
			name: "simple dockerfile",
			content: `FROM node:18-alpine
WORKDIR /app
COPY . .
RUN npm install
EXPOSE 3000
CMD ["node", "index.js"]`,
			expectedPorts: []int{3000},
			expectedArgs:  []string{},
			expectedBase:  "node:18-alpine",
		},
		{
			name: "multiple ports",
			content: `FROM golang:1.21
EXPOSE 8080
EXPOSE 9090
CMD ["./app"]`,
			expectedPorts: []int{8080, 9090},
			expectedArgs:  []string{},
			expectedBase:  "golang:1.21",
		},
		{
			name: "with build args",
			content: `FROM node:18
ARG NODE_ENV=production
ARG VERSION
ARG BUILD_DATE
EXPOSE 3000`,
			expectedPorts: []int{3000},
			expectedArgs:  []string{"NODE_ENV", "VERSION", "BUILD_DATE"},
			expectedBase:  "node:18",
		},
		{
			name: "multi-stage build",
			content: `FROM node:18 AS builder
WORKDIR /build
COPY . .
RUN npm run build

FROM node:18-alpine AS runtime
WORKDIR /app
COPY --from=builder /build/dist .
EXPOSE 8080`,
			expectedPorts: []int{8080},
			expectedArgs:  []string{},
			expectedBase:  "node:18", // First FROM is captured
		},
		{
			name: "with comments and empty lines",
			content: `# Base image
FROM python:3.11-slim

# Build arguments
ARG APP_VERSION
ARG DEBUG=false

# Expose port
EXPOSE 5000

# Run command
CMD ["python", "app.py"]`,
			expectedPorts: []int{5000},
			expectedArgs:  []string{"APP_VERSION", "DEBUG"},
			expectedBase:  "python:3.11-slim",
		},
		{
			name: "no expose or args",
			content: `FROM nginx:alpine
COPY nginx.conf /etc/nginx/nginx.conf
CMD ["nginx", "-g", "daemon off;"]`,
			expectedPorts: []int{},
			expectedArgs:  []string{},
			expectedBase:  "nginx:alpine",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &DockerfileInfo{
				ExposedPorts: []int{},
				BuildArgs:    []string{},
			}

			parseDockerfile(tt.content, info)

			assert.Equal(t, tt.expectedPorts, info.ExposedPorts, "exposed ports mismatch")
			assert.Equal(t, tt.expectedArgs, info.BuildArgs, "build args mismatch")
			assert.Equal(t, tt.expectedBase, info.BaseImage, "base image mismatch")
		})
	}
}

func TestParseDockerfile_MultiStage(t *testing.T) {
	content := `FROM golang:1.21 as builder
WORKDIR /build
COPY . .
RUN go build -o app

FROM alpine:latest
COPY --from=builder /build/app /app
EXPOSE 8080`

	info := &DockerfileInfo{
		ExposedPorts: []int{},
		BuildArgs:    []string{},
	}

	parseDockerfile(content, info)

	// Should capture first FROM, removing "as builder"
	assert.Equal(t, "golang:1.21", info.BaseImage)
	assert.Equal(t, []int{8080}, info.ExposedPorts)
}
