package store

import (
	"time"

	"gorm.io/datatypes"
)

type EnvironmentType string

const (
	EnvironmentTypeDev       EnvironmentType = "dev"
	EnvironmentTypeStaging   EnvironmentType = "staging"
	EnvironmentTypeProd      EnvironmentType = "prod"
	EnvironmentTypeEphemeral EnvironmentType = "ephemeral"
)

type Environment struct {
	ID                   string          `gorm:"primaryKey;type:text" json:"id"`
	Name                 string          `gorm:"index;not null" json:"name"`
	Type                 EnvironmentType `gorm:"index;not null" json:"type"`
	SourceRepo           string          `gorm:"type:text" json:"sourceRepo"`
	SourceBranch         string          `gorm:"type:text" json:"sourceBranch"`
	SourceCommit         string          `gorm:"type:text" json:"sourceCommit"`
	Status               string          `gorm:"index" json:"status"`
	RailwayProjectID     string          `gorm:"type:text" json:"railwayProjectId"` // Railway project ID (needed for provision outputs)
	RailwayEnvironmentID string          `gorm:"type:text" json:"railwayEnvironmentId"`
	TTLSeconds           *int64          `gorm:"type:integer" json:"ttlSeconds,omitempty"`
	ParentEnvironmentID  *string         `gorm:"type:text" json:"parentEnvironmentId,omitempty"`
	CreatedAt            time.Time       `gorm:"index" json:"createdAt"`
	UpdatedAt            time.Time       `json:"updatedAt"`

	Services []Service `gorm:"foreignKey:EnvironmentID" json:"services,omitempty"`
}

type DeploymentType string

const (
	DeploymentTypeSourceRepo  DeploymentType = "source_repo"
	DeploymentTypeDockerImage DeploymentType = "docker_image"
)

type Service struct {
	ID               string    `gorm:"primaryKey;type:text" json:"id"`
	EnvironmentID    string    `gorm:"index;not null" json:"environmentId"`
	Name             string    `gorm:"index;not null" json:"name"`
	Path             string    `gorm:"type:text" json:"path"`
	Status           string    `gorm:"index" json:"status"`
	RailwayServiceID string    `gorm:"type:text" json:"railwayServiceId"`
	CreatedAt        time.Time `gorm:"index" json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	// Deployment configuration
	DeploymentType DeploymentType `gorm:"type:text;default:'source_repo'" json:"deploymentType"`

	// Source repository fields (for source_repo deployment)
	SourceRepo   string `gorm:"type:text" json:"sourceRepo"`
	SourceBranch string `gorm:"type:text" json:"sourceBranch"`

	// Docker build configuration (for source_repo with Dockerfile)
	DockerfilePath *string `gorm:"type:text" json:"dockerfilePath,omitempty"` // Path to Dockerfile relative to repo root
	BuildContext   *string `gorm:"type:text" json:"buildContext,omitempty"`   // Docker build context path
	RootDirectory  *string `gorm:"type:text" json:"rootDirectory,omitempty"`  // Root directory for the service
	TargetStage    *string `gorm:"type:text" json:"targetStage,omitempty"`    // Multi-stage build target

	// Docker image fields (for docker_image deployment)
	DockerImage     string `gorm:"type:text" json:"dockerImage"`
	ImageRegistry   string `gorm:"type:text" json:"imageRegistry"`
	ImageName       string `gorm:"type:text" json:"imageName"`
	ImageTag        string `gorm:"type:text" json:"imageTag"`
	ImageDigest     string `gorm:"type:text" json:"imageDigest"`
	ImageAuthStored bool   `gorm:"default:false" json:"imageAuthStored"` // Indicates if Railway has stored auth credentials

	// Runtime configuration
	ExposedPortsJSON string  `gorm:"type:text" json:"exposedPortsJson"`          // JSON array of port numbers
	HealthCheckPath  *string `gorm:"type:text" json:"healthCheckPath,omitempty"` // Health check endpoint path
	StartCommand     *string `gorm:"type:text" json:"startCommand,omitempty"`    // Custom start command
}

// EnvironmentMetadata stores complete wizard state and provision outputs
// to enable environment cloning, branch-based deployments, and template creation.
type EnvironmentMetadata struct {
	ID            string `gorm:"primaryKey;type:text"`
	EnvironmentID string `gorm:"index;not null"` // Foreign key to Environment
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Template functionality
	IsTemplate          bool    `gorm:"index;default:false"` // Is this environment a template?
	TemplateName        *string `gorm:"type:text"`           // Template name (if IsTemplate=true)
	TemplateDescription *string `gorm:"type:text"`           // Template description

	// Cloning lineage
	ClonedFromEnvID *string `gorm:"type:text"` // ID of environment this was cloned from

	// Wizard state and provision outputs (stored as JSON for flexibility)
	WizardInputsJSON     datatypes.JSON `gorm:"type:jsonb"` // Complete wizard state (all inputs from all steps)
	ProvisionOutputsJSON datatypes.JSON `gorm:"type:jsonb"` // Provision outputs (Railway project/environment/service IDs)
}
