package store

import (
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Open opens a GORM DB with SQLite and runs AutoMigrate for models.
func Open(sqlitePath string) (*gorm.DB, error) {
	if sqlitePath == "" {
		sqlitePath = "mirage.db"
	}
	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&User{}, &Environment{}, &Service{}, &EnvironmentMetadata{}); err != nil {
		return nil, err
	}
	return db, nil
}

// OpenFromURL opens a GORM DB based on a database URL. Supports Postgres URLs
// (postgres:// or postgresql://). Falls back to SQLite for empty or non-PG URLs.
func OpenFromURL(databaseURL string) (*gorm.DB, error) {
	trimmed := strings.TrimSpace(databaseURL)
	if trimmed == "" {
		// default to local sqlite file
		return Open("")
	}
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "postgres://") || strings.HasPrefix(lower, "postgresql://") {
		db, err := gorm.Open(postgres.Open(trimmed), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		if err := db.AutoMigrate(&User{}, &Environment{}, &Service{}, &EnvironmentMetadata{}); err != nil {
			return nil, err
		}
		return db, nil
	}
	// Otherwise treat as sqlite DSN or path
	return Open(trimmed)
}
