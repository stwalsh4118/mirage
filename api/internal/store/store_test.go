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

	// Marshal ports to JSON
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

	// Verify ExposedPortsJSON can be unmarshaled
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

func TestEnvironmentMetadata_CreateWithWizardInputs(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	// Create a parent environment first
	env := Environment{
		ID:        "env-1",
		Name:      "test-env",
		Type:      EnvironmentTypeDev,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&env).Error; err != nil {
		t.Fatalf("failed to create environment: %v", err)
	}

	// Create complex wizard inputs JSON
	wizardInputs := map[string]interface{}{
		"step1": map[string]interface{}{
			"projectName": "my-project",
			"repository":  "github.com/owner/repo",
			"branch":      "main",
		},
		"step2": map[string]interface{}{
			"services": []map[string]interface{}{
				{
					"name":           "api",
					"path":           "services/api",
					"dockerfilePath": "services/api/Dockerfile",
					"port":           3000,
				},
				{
					"name": "worker",
					"path": "services/worker",
				},
			},
		},
		"step3": map[string]interface{}{
			"environmentType": "dev",
			"variables": map[string]string{
				"NODE_ENV":     "development",
				"DATABASE_URL": "postgres://localhost:5432/mydb",
			},
		},
	}
	wizardInputsJSON, _ := json.Marshal(wizardInputs)

	// Create provision outputs JSON
	provisionOutputs := map[string]interface{}{
		"projectId":     "railway-proj-123",
		"environmentId": "railway-env-456",
		"services": []map[string]string{
			{"name": "api", "serviceId": "railway-svc-789"},
			{"name": "worker", "serviceId": "railway-svc-abc"},
		},
	}
	provisionOutputsJSON, _ := json.Marshal(provisionOutputs)

	metadata := EnvironmentMetadata{
		ID:                   "meta-1",
		EnvironmentID:        "env-1",
		IsTemplate:           false,
		WizardInputsJSON:     wizardInputsJSON,
		ProvisionOutputsJSON: provisionOutputsJSON,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := db.Create(&metadata).Error; err != nil {
		t.Fatalf("failed to create environment metadata: %v", err)
	}

	var retrieved EnvironmentMetadata
	if err := db.First(&retrieved, "id = ?", "meta-1").Error; err != nil {
		t.Fatalf("failed to retrieve metadata: %v", err)
	}

	// Verify foreign key
	if retrieved.EnvironmentID != "env-1" {
		t.Errorf("expected environment ID env-1, got %s", retrieved.EnvironmentID)
	}

	// Verify we can unmarshal wizard inputs
	var retrievedWizardInputs map[string]interface{}
	if err := json.Unmarshal(retrieved.WizardInputsJSON, &retrievedWizardInputs); err != nil {
		t.Fatalf("failed to unmarshal wizard inputs: %v", err)
	}

	step1 := retrievedWizardInputs["step1"].(map[string]interface{})
	if step1["projectName"] != "my-project" {
		t.Errorf("expected project name 'my-project', got %v", step1["projectName"])
	}

	// Verify we can unmarshal provision outputs
	var retrievedProvisionOutputs map[string]interface{}
	if err := json.Unmarshal(retrieved.ProvisionOutputsJSON, &retrievedProvisionOutputs); err != nil {
		t.Fatalf("failed to unmarshal provision outputs: %v", err)
	}

	if retrievedProvisionOutputs["projectId"] != "railway-proj-123" {
		t.Errorf("expected project ID 'railway-proj-123', got %v", retrievedProvisionOutputs["projectId"])
	}

	// Verify timestamps
	if retrieved.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt timestamp")
	}
	if retrieved.UpdatedAt.IsZero() {
		t.Error("expected non-zero UpdatedAt timestamp")
	}
}

func TestEnvironmentMetadata_QueryByEnvironmentID(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	// Create parent environment
	env := Environment{
		ID:        "env-1",
		Name:      "test-env",
		Type:      EnvironmentTypeDev,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&env).Error; err != nil {
		t.Fatalf("failed to create environment: %v", err)
	}

	// Create metadata
	metadata := EnvironmentMetadata{
		ID:                   "meta-1",
		EnvironmentID:        "env-1",
		WizardInputsJSON:     []byte("{}"),
		ProvisionOutputsJSON: []byte("{}"),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	if err := db.Create(&metadata).Error; err != nil {
		t.Fatalf("failed to create metadata: %v", err)
	}

	// Query by environment ID
	var retrieved EnvironmentMetadata
	if err := db.Where("environment_id = ?", "env-1").First(&retrieved).Error; err != nil {
		t.Fatalf("failed to query by environment ID: %v", err)
	}

	if retrieved.ID != "meta-1" {
		t.Errorf("expected metadata ID meta-1, got %s", retrieved.ID)
	}
}

func TestEnvironmentMetadata_TemplateFields(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	// Create parent environment
	env := Environment{
		ID:        "env-template",
		Name:      "template-env",
		Type:      EnvironmentTypeDev,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&env).Error; err != nil {
		t.Fatalf("failed to create environment: %v", err)
	}

	templateName := "Node.js Microservices Stack"
	templateDesc := "Complete microservices setup with API, worker, and database"

	metadata := EnvironmentMetadata{
		ID:                   "meta-template",
		EnvironmentID:        "env-template",
		IsTemplate:           true,
		TemplateName:         &templateName,
		TemplateDescription:  &templateDesc,
		WizardInputsJSON:     []byte("{}"),
		ProvisionOutputsJSON: []byte("{}"),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := db.Create(&metadata).Error; err != nil {
		t.Fatalf("failed to create template metadata: %v", err)
	}

	// Query templates
	var templates []EnvironmentMetadata
	if err := db.Where("is_template = ?", true).Find(&templates).Error; err != nil {
		t.Fatalf("failed to query templates: %v", err)
	}

	if len(templates) != 1 {
		t.Fatalf("expected 1 template, got %d", len(templates))
	}

	retrieved := templates[0]
	if !retrieved.IsTemplate {
		t.Error("expected IsTemplate to be true")
	}
	if retrieved.TemplateName == nil || *retrieved.TemplateName != templateName {
		t.Errorf("expected template name '%s', got %v", templateName, retrieved.TemplateName)
	}
	if retrieved.TemplateDescription == nil || *retrieved.TemplateDescription != templateDesc {
		t.Errorf("expected template description '%s', got %v", templateDesc, retrieved.TemplateDescription)
	}
}

func TestEnvironmentMetadata_ClonedFromEnvID(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	// Create original environment
	origEnv := Environment{
		ID:        "env-original",
		Name:      "original-env",
		Type:      EnvironmentTypeDev,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&origEnv).Error; err != nil {
		t.Fatalf("failed to create original environment: %v", err)
	}

	// Create cloned environment
	clonedEnv := Environment{
		ID:        "env-cloned",
		Name:      "cloned-env",
		Type:      EnvironmentTypeStaging,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&clonedEnv).Error; err != nil {
		t.Fatalf("failed to create cloned environment: %v", err)
	}

	clonedFromID := "env-original"
	metadata := EnvironmentMetadata{
		ID:                   "meta-cloned",
		EnvironmentID:        "env-cloned",
		ClonedFromEnvID:      &clonedFromID,
		WizardInputsJSON:     []byte("{}"),
		ProvisionOutputsJSON: []byte("{}"),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := db.Create(&metadata).Error; err != nil {
		t.Fatalf("failed to create cloned metadata: %v", err)
	}

	var retrieved EnvironmentMetadata
	if err := db.First(&retrieved, "id = ?", "meta-cloned").Error; err != nil {
		t.Fatalf("failed to retrieve metadata: %v", err)
	}

	if retrieved.ClonedFromEnvID == nil || *retrieved.ClonedFromEnvID != "env-original" {
		t.Errorf("expected cloned from env ID 'env-original', got %v", retrieved.ClonedFromEnvID)
	}
}

func TestEnvironmentMetadata_OptionalFieldsNil(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	// Create parent environment
	env := Environment{
		ID:        "env-1",
		Name:      "minimal-env",
		Type:      EnvironmentTypeDev,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&env).Error; err != nil {
		t.Fatalf("failed to create environment: %v", err)
	}

	// Create metadata with minimal fields
	metadata := EnvironmentMetadata{
		ID:                   "meta-minimal",
		EnvironmentID:        "env-1",
		IsTemplate:           false,
		WizardInputsJSON:     []byte("{}"),
		ProvisionOutputsJSON: []byte("{}"),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := db.Create(&metadata).Error; err != nil {
		t.Fatalf("failed to create minimal metadata: %v", err)
	}

	var retrieved EnvironmentMetadata
	if err := db.First(&retrieved, "id = ?", "meta-minimal").Error; err != nil {
		t.Fatalf("failed to retrieve metadata: %v", err)
	}

	// Verify optional fields are nil
	if retrieved.TemplateName != nil {
		t.Errorf("expected nil TemplateName, got %v", *retrieved.TemplateName)
	}
	if retrieved.TemplateDescription != nil {
		t.Errorf("expected nil TemplateDescription, got %v", *retrieved.TemplateDescription)
	}
	if retrieved.ClonedFromEnvID != nil {
		t.Errorf("expected nil ClonedFromEnvID, got %v", *retrieved.ClonedFromEnvID)
	}
	if retrieved.IsTemplate {
		t.Error("expected IsTemplate to be false")
	}
}

func TestEnvironmentMetadata_ComplexNestedJSON(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	// Create parent environment
	env := Environment{
		ID:        "env-1",
		Name:      "complex-env",
		Type:      EnvironmentTypeDev,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&env).Error; err != nil {
		t.Fatalf("failed to create environment: %v", err)
	}

	// Create deeply nested complex structure
	complexInputs := map[string]interface{}{
		"configuration": map[string]interface{}{
			"networking": map[string]interface{}{
				"ports": []int{3000, 8080, 9090},
				"hostnames": []string{
					"api.example.com",
					"worker.example.com",
				},
				"ssl": map[string]bool{
					"enabled":    true,
					"autoRenew":  true,
					"forceHTTPS": true,
				},
			},
			"scaling": map[string]interface{}{
				"min": 1,
				"max": 10,
				"targets": map[string]float64{
					"cpu":    0.7,
					"memory": 0.8,
				},
			},
		},
		"services": []map[string]interface{}{
			{
				"name": "api",
				"config": map[string]interface{}{
					"buildArgs": map[string]string{
						"NODE_ENV": "production",
						"VERSION":  "1.2.3",
					},
					"healthCheck": map[string]interface{}{
						"path":     "/health",
						"interval": 30,
						"timeout":  10,
					},
				},
			},
		},
	}
	complexInputsJSON, _ := json.Marshal(complexInputs)

	metadata := EnvironmentMetadata{
		ID:                   "meta-complex",
		EnvironmentID:        "env-1",
		WizardInputsJSON:     complexInputsJSON,
		ProvisionOutputsJSON: []byte("{}"),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := db.Create(&metadata).Error; err != nil {
		t.Fatalf("failed to create metadata with complex JSON: %v", err)
	}

	var retrieved EnvironmentMetadata
	if err := db.First(&retrieved, "id = ?", "meta-complex").Error; err != nil {
		t.Fatalf("failed to retrieve metadata: %v", err)
	}

	// Verify complex nested structure is preserved
	var retrievedInputs map[string]interface{}
	if err := json.Unmarshal(retrieved.WizardInputsJSON, &retrievedInputs); err != nil {
		t.Fatalf("failed to unmarshal complex inputs: %v", err)
	}

	config := retrievedInputs["configuration"].(map[string]interface{})
	networking := config["networking"].(map[string]interface{})
	ports := networking["ports"].([]interface{})

	if len(ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(ports))
	}

	ssl := networking["ssl"].(map[string]interface{})
	if ssl["enabled"] != true {
		t.Error("expected ssl.enabled to be true")
	}

	services := retrievedInputs["services"].([]interface{})
	if len(services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(services))
	}

	service := services[0].(map[string]interface{})
	serviceConfig := service["config"].(map[string]interface{})
	buildArgs := serviceConfig["buildArgs"].(map[string]interface{})

	if buildArgs["NODE_ENV"] != "production" {
		t.Errorf("expected NODE_ENV=production, got %v", buildArgs["NODE_ENV"])
	}
}
