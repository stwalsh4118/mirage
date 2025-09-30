package store

import (
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
