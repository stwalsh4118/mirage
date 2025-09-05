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
	RailwayEnvironmentID string          `gorm:"type:text"`
	TTLSeconds           *int64          `gorm:"type:integer"`
	ParentEnvironmentID  *string         `gorm:"type:text"`
	CreatedAt            time.Time       `gorm:"index"`
	UpdatedAt            time.Time

	Services []Service `gorm:"foreignKey:EnvironmentID"`
}

type Service struct {
	ID               string    `gorm:"primaryKey;type:text"`
	EnvironmentID    string    `gorm:"index;not null"`
	Name             string    `gorm:"index;not null"`
	Path             string    `gorm:"type:text"`
	Status           string    `gorm:"index"`
	RailwayServiceID string    `gorm:"type:text"`
	CreatedAt        time.Time `gorm:"index"`
	UpdatedAt        time.Time
}

type Template struct {
	ID        string    `gorm:"primaryKey;type:text"`
	Name      string    `gorm:"uniqueIndex;not null"`
	BaseJSON  string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}
