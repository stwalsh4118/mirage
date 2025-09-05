package store

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Open opens a GORM DB with SQLite and runs AutoMigrate for models.
func Open(sqlitePath string) (*gorm.DB, error) {
	if sqlitePath == "" {
		sqlitePath = "mirage.db"
	}
	dsn := fmt.Sprintf("%s", sqlitePath)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Environment{}, &Service{}, &Template{}); err != nil {
		return nil, err
	}
	return db, nil
}
