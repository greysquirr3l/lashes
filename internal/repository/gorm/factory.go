package gorm

import (
	"fmt"
	"time"

	"github.com/greysquirr3l/lashes/internal/storage"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewDB creates a new GORM database connection using the provided options
func NewDB(opts storage.Options) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch opts.Type {
	case storage.SQLite:
		dialector = sqlite.Open(opts.FilePath)
	case storage.MySQL:
		dialector = mysql.Open(opts.ConnectionString)
	case storage.Postgres:
		dialector = postgres.Open(opts.ConnectionString)
	case storage.Memory:
		// Memory is not supported by GORM, fall through to default case
		fallthrough
	default:
		return nil, fmt.Errorf("unsupported database type: %s", opts.Type)
	}

	config := &gorm.Config{}

	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}

	// Configure the connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set reasonable defaults
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto migrate models
	if err := db.AutoMigrate(&ProxyModel{}); err != nil {
		return nil, err
	}

	return db, nil
}
