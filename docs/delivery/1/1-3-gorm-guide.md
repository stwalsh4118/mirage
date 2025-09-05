# 1-3 GORM Guide (2025-09-05)

## Packages
- ORM: `gorm.io/gorm`
- SQLite driver (MVP): `gorm.io/driver/sqlite`

## Install
```bash
go get gorm.io/gorm gorm.io/driver/sqlite
```

## Initialize (SQLite)
```go
import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

db, err := gorm.Open(sqlite.Open("mirage.db"), &gorm.Config{})
```

## AutoMigrate
```go
if err := db.AutoMigrate(&Environment{}, &Service{}, &Template{}); err != nil {
	// handle
}
```

## References
- GORM Docs: `https://gorm.io/docs/`
- SQLite Driver: `https://gorm.io/docs/connecting_to_the_database.html#SQLite`
