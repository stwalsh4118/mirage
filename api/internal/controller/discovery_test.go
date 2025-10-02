package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stwalsh4118/mirageapi/internal/scanner"
)

// mockDockerfileScanner implements DockerfileScanner for testing
type mockDockerfileScanner struct {
	scanFunc func(ctx context.Context, owner, repo, branch, userToken string) ([]scanner.DockerfileInfo, error)
}

func (m *mockDockerfileScanner) ScanRepository(ctx context.Context, owner, repo, branch, userToken string) ([]scanner.DockerfileInfo, error) {
	if m.scanFunc != nil {
		return m.scanFunc(ctx, owner, repo, branch, userToken)
	}
	return []scanner.DockerfileInfo{}, nil
}

func TestDiscoverDockerfiles_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockScanner := &mockDockerfileScanner{
		scanFunc: func(ctx context.Context, owner, repo, branch, userToken string) ([]scanner.DockerfileInfo, error) {
			assert.Equal(t, "test-owner", owner)
			assert.Equal(t, "test-repo", repo)
			assert.Equal(t, "main", branch)

			return []scanner.DockerfileInfo{
				{
					Path:         "services/api/Dockerfile",
					ServiceName:  "api",
					BuildContext: "services/api",
					ExposedPorts: []int{3000, 9090},
					BuildArgs:    []string{"NODE_ENV", "VERSION"},
					BaseImage:    "node:18-alpine",
				},
				{
					Path:         "services/worker/Dockerfile",
					ServiceName:  "worker",
					BuildContext: "services/worker",
					ExposedPorts: []int{},
					BuildArgs:    []string{},
					BaseImage:    "node:18-alpine",
				},
			}, nil
		},
	}

	controller := &DiscoveryController{Scanner: mockScanner}

	reqBody := DiscoverDockerfilesRequest{
		Owner:     "test-owner",
		Repo:      "test-repo",
		Branch:    "main",
		RequestID: "test-123",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/discovery/dockerfiles", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.DiscoverDockerfiles(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp DiscoverDockerfilesResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, "test-owner", resp.Owner)
	assert.Equal(t, "test-repo", resp.Repo)
	assert.Equal(t, "main", resp.Branch)
	assert.Len(t, resp.Services, 2)

	// Verify first service
	assert.Equal(t, "api", resp.Services[0].Name)
	assert.Equal(t, "services/api/Dockerfile", resp.Services[0].DockerfilePath)
	assert.Equal(t, "services/api", resp.Services[0].BuildContext)
	assert.Equal(t, []int{3000, 9090}, resp.Services[0].ExposedPorts)
	assert.Equal(t, []string{"NODE_ENV", "VERSION"}, resp.Services[0].BuildArgs)
	assert.Equal(t, "node:18-alpine", resp.Services[0].BaseImage)

	// Verify second service
	assert.Equal(t, "worker", resp.Services[1].Name)
	assert.Equal(t, "services/worker/Dockerfile", resp.Services[1].DockerfilePath)
}

func TestDiscoverDockerfiles_NoDockerfilesFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockScanner := &mockDockerfileScanner{
		scanFunc: func(ctx context.Context, owner, repo, branch, userToken string) ([]scanner.DockerfileInfo, error) {
			return []scanner.DockerfileInfo{}, nil
		},
	}

	controller := &DiscoveryController{Scanner: mockScanner}

	reqBody := DiscoverDockerfilesRequest{
		Owner:  "test-owner",
		Repo:   "empty-repo",
		Branch: "main",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/discovery/dockerfiles", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.DiscoverDockerfiles(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp DiscoverDockerfilesResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Len(t, resp.Services, 0)
}

func TestDiscoverDockerfiles_ScannerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockScanner := &mockDockerfileScanner{
		scanFunc: func(ctx context.Context, owner, repo, branch, userToken string) ([]scanner.DockerfileInfo, error) {
			return nil, fmt.Errorf("repository not found: test-owner/invalid-repo@main")
		},
	}

	controller := &DiscoveryController{Scanner: mockScanner}

	reqBody := DiscoverDockerfilesRequest{
		Owner:  "test-owner",
		Repo:   "invalid-repo",
		Branch: "main",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/discovery/dockerfiles", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.DiscoverDockerfiles(c)

	assert.Equal(t, http.StatusBadGateway, w.Code)
	assert.Contains(t, w.Body.String(), "failed to scan repository")
}

func TestDiscoverDockerfiles_WithUserToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var capturedToken string
	mockScanner := &mockDockerfileScanner{
		scanFunc: func(ctx context.Context, owner, repo, branch, userToken string) ([]scanner.DockerfileInfo, error) {
			capturedToken = userToken
			return []scanner.DockerfileInfo{}, nil
		},
	}

	controller := &DiscoveryController{Scanner: mockScanner}

	reqBody := DiscoverDockerfilesRequest{
		Owner:     "test-owner",
		Repo:      "private-repo",
		Branch:    "main",
		UserToken: "ghp_test_token_123",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/discovery/dockerfiles", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.DiscoverDockerfiles(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ghp_test_token_123", capturedToken)
}

func TestDiscoverDockerfiles_ValidationErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		request  DiscoverDockerfilesRequest
		expected string
	}{
		{
			name: "missing owner",
			request: DiscoverDockerfilesRequest{
				Repo:   "test-repo",
				Branch: "main",
			},
			expected: "Owner",
		},
		{
			name: "missing repo",
			request: DiscoverDockerfilesRequest{
				Owner:  "test-owner",
				Branch: "main",
			},
			expected: "Repo",
		},
		{
			name: "missing branch",
			request: DiscoverDockerfilesRequest{
				Owner: "test-owner",
				Repo:  "test-repo",
			},
			expected: "Branch",
		},
		{
			name: "invalid owner format",
			request: DiscoverDockerfilesRequest{
				Owner:  "test@owner",
				Repo:   "test-repo",
				Branch: "main",
			},
			expected: "invalid owner format",
		},
		{
			name: "invalid repo format",
			request: DiscoverDockerfilesRequest{
				Owner:  "test-owner",
				Repo:   "test@repo",
				Branch: "main",
			},
			expected: "invalid repo format",
		},
		{
			name: "invalid branch format",
			request: DiscoverDockerfilesRequest{
				Owner:  "test-owner",
				Repo:   "test-repo",
				Branch: ".invalid",
			},
			expected: "invalid branch format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := &DiscoveryController{Scanner: &mockDockerfileScanner{}}

			body, _ := json.Marshal(tt.request)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/discovery/dockerfiles", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			controller.DiscoverDockerfiles(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), tt.expected)
		})
	}
}

func TestDiscoverDockerfiles_ScannerNotConfigured(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := &DiscoveryController{Scanner: nil}

	reqBody := DiscoverDockerfilesRequest{
		Owner:  "test-owner",
		Repo:   "test-repo",
		Branch: "main",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/discovery/dockerfiles", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.DiscoverDockerfiles(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), "scanner not configured")
}

func TestValidateGitHubIdentifier(t *testing.T) {
	validIdentifiers := []string{
		"user",
		"test-user",
		"test_user",
		"test.user",
		"User123",
		"org-name",
	}

	for _, id := range validIdentifiers {
		t.Run(id, func(t *testing.T) {
			valid := isValidGitHubIdentifier(id)
			assert.True(t, valid, "identifier %q should be valid", id)
		})
	}

	invalidIdentifiers := []string{
		"",                        // empty
		"user@domain",             // @ not allowed
		"user/repo",               // / not allowed
		"user:name",               // : not allowed
		"user name",               // space not allowed
		"user#123",                // # not allowed
		string(make([]byte, 101)), // too long
	}

	for _, id := range invalidIdentifiers {
		t.Run(fmt.Sprintf("invalid_%s", id), func(t *testing.T) {
			valid := isValidGitHubIdentifier(id)
			assert.False(t, valid, "identifier %q should be invalid", id)
		})
	}
}

func TestValidateBranchName(t *testing.T) {
	validBranches := []string{
		"main",
		"master",
		"develop",
		"feature/my-feature",
		"release-1.0",
		"bugfix_123",
		"v1.0.0",
	}

	for _, branch := range validBranches {
		t.Run(branch, func(t *testing.T) {
			valid := isValidBranchName(branch)
			assert.True(t, valid, "branch %q should be valid", branch)
		})
	}

	invalidBranches := []string{
		"",                        // empty
		".main",                   // starts with dot
		"/main",                   // starts with slash
		"branch name",             // space not allowed
		"branch@123",              // @ not allowed
		string(make([]byte, 256)), // too long
	}

	for _, branch := range invalidBranches {
		t.Run(fmt.Sprintf("invalid_%s", branch), func(t *testing.T) {
			valid := isValidBranchName(branch)
			assert.False(t, valid, "branch %q should be invalid", branch)
		})
	}
}
