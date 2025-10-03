package store

import (
	"encoding/json"
	"testing"
	"time"
)

func TestOpen_InMemorySQLite_AutoMigrate(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	if db == nil {
		t.Fatalf("expected db instance, got nil")
	}
}

func TestOpenFromURL_FallbackToSQLite(t *testing.T) {
	db, err := OpenFromURL("")
	if err != nil {
		t.Fatalf("open from url failed: %v", err)
	}
	if db == nil {
		t.Fatalf("expected db instance, got nil")
	}
}

func TestService_SourceRepoDeployment(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	service := Service{
		ID:             "svc-1",
		EnvironmentID:  "env-1",
		Name:           "api",
		DeploymentType: DeploymentTypeSourceRepo,
		SourceRepo:     "github.com/owner/repo",
		SourceBranch:   "main",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := db.Create(&service).Error; err != nil {
		t.Fatalf("failed to create source repo service: %v", err)
	}

	var retrieved Service
	if err := db.First(&retrieved, "id = ?", "svc-1").Error; err != nil {
		t.Fatalf("failed to retrieve service: %v", err)
	}

	if retrieved.DeploymentType != DeploymentTypeSourceRepo {
		t.Errorf("expected deployment type %s, got %s", DeploymentTypeSourceRepo, retrieved.DeploymentType)
	}
	if retrieved.SourceRepo != "github.com/owner/repo" {
		t.Errorf("expected source repo github.com/owner/repo, got %s", retrieved.SourceRepo)
	}
}

func TestService_DockerImageDeployment(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	service := Service{
		ID:              "svc-2",
		EnvironmentID:   "env-1",
		Name:            "nginx",
		DeploymentType:  DeploymentTypeDockerImage,
		DockerImage:     "nginx:latest",
		ImageRegistry:   "docker.io",
		ImageName:       "nginx",
		ImageTag:        "latest",
		ImageAuthStored: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := db.Create(&service).Error; err != nil {
		t.Fatalf("failed to create docker image service: %v", err)
	}

	var retrieved Service
	if err := db.First(&retrieved, "id = ?", "svc-2").Error; err != nil {
		t.Fatalf("failed to retrieve service: %v", err)
	}

	if retrieved.DeploymentType != DeploymentTypeDockerImage {
		t.Errorf("expected deployment type %s, got %s", DeploymentTypeDockerImage, retrieved.DeploymentType)
	}
	if retrieved.DockerImage != "nginx:latest" {
		t.Errorf("expected docker image nginx:latest, got %s", retrieved.DockerImage)
	}
	if retrieved.ImageRegistry != "docker.io" {
		t.Errorf("expected registry docker.io, got %s", retrieved.ImageRegistry)
	}
}

func TestService_DefaultDeploymentType(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	// Create service without specifying deployment type
	service := Service{
		ID:            "svc-3",
		EnvironmentID: "env-1",
		Name:          "legacy-service",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := db.Create(&service).Error; err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	var retrieved Service
	if err := db.First(&retrieved, "id = ?", "svc-3").Error; err != nil {
		t.Fatalf("failed to retrieve service: %v", err)
	}

	// Should default to source_repo
	if retrieved.DeploymentType != DeploymentTypeSourceRepo {
		t.Errorf("expected default deployment type %s, got %s", DeploymentTypeSourceRepo, retrieved.DeploymentType)
	}
}

func TestService_BuildConfiguration(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	dockerfilePath := "services/api/Dockerfile"
	buildContext := "services/api"
	rootDir := "services/api"
	healthCheck := "/health"
	startCmd := "npm start"
	targetStage := "production"

	// Marshal build args and ports to JSON
	buildArgs := map[string]string{"NODE_ENV": "production", "VERSION": "1.0.0"}
	buildArgsJSON, _ := json.Marshal(buildArgs)
	exposedPorts := []int{3000, 8080}
	exposedPortsJSON, _ := json.Marshal(exposedPorts)

	service := Service{
		ID:               "svc-build",
		EnvironmentID:    "env-1",
		Name:             "api-service",
		DeploymentType:   DeploymentTypeSourceRepo,
		SourceRepo:       "github.com/owner/monorepo",
		SourceBranch:     "main",
		DockerfilePath:   &dockerfilePath,
		BuildContext:     &buildContext,
		RootDirectory:    &rootDir,
		BuildArgsJSON:    string(buildArgsJSON),
		TargetStage:      &targetStage,
		ExposedPortsJSON: string(exposedPortsJSON),
		HealthCheckPath:  &healthCheck,
		StartCommand:     &startCmd,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := db.Create(&service).Error; err != nil {
		t.Fatalf("failed to create service with build config: %v", err)
	}

	var retrieved Service
	if err := db.First(&retrieved, "id = ?", "svc-build").Error; err != nil {
		t.Fatalf("failed to retrieve service: %v", err)
	}

	// Verify all build configuration fields
	if retrieved.DockerfilePath == nil || *retrieved.DockerfilePath != dockerfilePath {
		t.Errorf("expected dockerfile path %s, got %v", dockerfilePath, retrieved.DockerfilePath)
	}
	if retrieved.BuildContext == nil || *retrieved.BuildContext != buildContext {
		t.Errorf("expected build context %s, got %v", buildContext, retrieved.BuildContext)
	}
	if retrieved.RootDirectory == nil || *retrieved.RootDirectory != rootDir {
		t.Errorf("expected root directory %s, got %v", rootDir, retrieved.RootDirectory)
	}
	if retrieved.TargetStage == nil || *retrieved.TargetStage != targetStage {
		t.Errorf("expected target stage %s, got %v", targetStage, retrieved.TargetStage)
	}
	if retrieved.HealthCheckPath == nil || *retrieved.HealthCheckPath != healthCheck {
		t.Errorf("expected health check path %s, got %v", healthCheck, retrieved.HealthCheckPath)
	}
	if retrieved.StartCommand == nil || *retrieved.StartCommand != startCmd {
		t.Errorf("expected start command %s, got %v", startCmd, retrieved.StartCommand)
	}

	// Verify JSON fields can be unmarshaled
	var retrievedBuildArgs map[string]string
	if err := json.Unmarshal([]byte(retrieved.BuildArgsJSON), &retrievedBuildArgs); err != nil {
		t.Fatalf("failed to unmarshal build args: %v", err)
	}
	if retrievedBuildArgs["NODE_ENV"] != "production" {
		t.Errorf("expected NODE_ENV=production in build args, got %s", retrievedBuildArgs["NODE_ENV"])
	}

	var retrievedPorts []int
	if err := json.Unmarshal([]byte(retrieved.ExposedPortsJSON), &retrievedPorts); err != nil {
		t.Fatalf("failed to unmarshal exposed ports: %v", err)
	}
	if len(retrievedPorts) != 2 || retrievedPorts[0] != 3000 {
		t.Errorf("expected ports [3000, 8080], got %v", retrievedPorts)
	}
}

func TestService_OptionalFieldsNil(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	// Create service with all optional fields as nil
	service := Service{
		ID:             "svc-minimal",
		EnvironmentID:  "env-1",
		Name:           "minimal-service",
		DeploymentType: DeploymentTypeSourceRepo,
		SourceRepo:     "github.com/owner/repo",
		SourceBranch:   "main",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := db.Create(&service).Error; err != nil {
		t.Fatalf("failed to create minimal service: %v", err)
	}

	var retrieved Service
	if err := db.First(&retrieved, "id = ?", "svc-minimal").Error; err != nil {
		t.Fatalf("failed to retrieve service: %v", err)
	}

	// Verify optional fields are nil
	if retrieved.DockerfilePath != nil {
		t.Errorf("expected nil DockerfilePath, got %v", *retrieved.DockerfilePath)
	}
	if retrieved.BuildContext != nil {
		t.Errorf("expected nil BuildContext, got %v", *retrieved.BuildContext)
	}
	if retrieved.RootDirectory != nil {
		t.Errorf("expected nil RootDirectory, got %v", *retrieved.RootDirectory)
	}
	if retrieved.TargetStage != nil {
		t.Errorf("expected nil TargetStage, got %v", *retrieved.TargetStage)
	}
	if retrieved.HealthCheckPath != nil {
		t.Errorf("expected nil HealthCheckPath, got %v", *retrieved.HealthCheckPath)
	}
	if retrieved.StartCommand != nil {
		t.Errorf("expected nil StartCommand, got %v", *retrieved.StartCommand)
	}
}

func TestEnvironment_RailwayProjectID(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	env := Environment{
		ID:                   "env-1",
		Name:                 "test-env",
		Type:                 EnvironmentTypeDev,
		RailwayProjectID:     "proj-123",
		RailwayEnvironmentID: "env-456",
		Status:               "active",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := db.Create(&env).Error; err != nil {
		t.Fatalf("failed to create environment: %v", err)
	}

	var retrieved Environment
	if err := db.First(&retrieved, "id = ?", "env-1").Error; err != nil {
		t.Fatalf("failed to retrieve environment: %v", err)
	}

	if retrieved.RailwayProjectID != "proj-123" {
		t.Errorf("expected railway project ID proj-123, got %s", retrieved.RailwayProjectID)
	}
	if retrieved.RailwayEnvironmentID != "env-456" {
		t.Errorf("expected railway environment ID env-456, got %s", retrieved.RailwayEnvironmentID)
	}
}
