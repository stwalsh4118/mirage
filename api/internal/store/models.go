package store

import "time"

type EnvironmentType string

const (
	EnvironmentTypeDev       EnvironmentType = "dev"
	EnvironmentTypeStaging   EnvironmentType = "staging"
	EnvironmentTypeProd      EnvironmentType = "prod"
	EnvironmentTypeEphemeral EnvironmentType = "ephemeral"
)

type Environment struct {
	ID                   string          `gorm:"primaryKey;type:text"`
	Name                 string          `gorm:"index;not null"`
	Type                 EnvironmentType `gorm:"index;not null"`
	SourceRepo           string          `gorm:"type:text"`
	SourceBranch         string          `gorm:"type:text"`
	SourceCommit         string          `gorm:"type:text"`
	Status               string          `gorm:"index"`
	RailwayProjectID     string          `gorm:"type:text"` // Railway project ID (needed for provision outputs)
	RailwayEnvironmentID string          `gorm:"type:text"`
	TTLSeconds           *int64          `gorm:"type:integer"`
	ParentEnvironmentID  *string         `gorm:"type:text"`
	CreatedAt            time.Time       `gorm:"index"`
	UpdatedAt            time.Time

	Services []Service `gorm:"foreignKey:EnvironmentID"`
}

type DeploymentType string

const (
	DeploymentTypeSourceRepo  DeploymentType = "source_repo"
	DeploymentTypeDockerImage DeploymentType = "docker_image"
)

type Service struct {
	ID               string    `gorm:"primaryKey;type:text"`
	EnvironmentID    string    `gorm:"index;not null"`
	Name             string    `gorm:"index;not null"`
	Path             string    `gorm:"type:text"`
	Status           string    `gorm:"index"`
	RailwayServiceID string    `gorm:"type:text"`
	CreatedAt        time.Time `gorm:"index"`
	UpdatedAt        time.Time

	// Deployment configuration
	DeploymentType DeploymentType `gorm:"type:text;default:'source_repo'"`

	// Source repository fields (for source_repo deployment)
	SourceRepo   string `gorm:"type:text"`
	SourceBranch string `gorm:"type:text"`

	// Docker build configuration (for source_repo with Dockerfile)
	DockerfilePath *string `gorm:"type:text"` // Path to Dockerfile relative to repo root
	BuildContext   *string `gorm:"type:text"` // Docker build context path
	RootDirectory  *string `gorm:"type:text"` // Root directory for the service
	BuildArgsJSON  string  `gorm:"type:text"` // JSON array/map of build arguments
	TargetStage    *string `gorm:"type:text"` // Multi-stage build target

	// Docker image fields (for docker_image deployment)
	DockerImage     string `gorm:"type:text"`
	ImageRegistry   string `gorm:"type:text"`
	ImageName       string `gorm:"type:text"`
	ImageTag        string `gorm:"type:text"`
	ImageDigest     string `gorm:"type:text"`
	ImageAuthStored bool   `gorm:"default:false"` // Indicates if Railway has stored auth credentials

	// Runtime configuration
	ExposedPortsJSON string  `gorm:"type:text"` // JSON array of port numbers
	HealthCheckPath  *string `gorm:"type:text"` // Health check endpoint path
	StartCommand     *string `gorm:"type:text"` // Custom start command
}
