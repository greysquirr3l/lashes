package gorm

import (
	"fmt"

	"github.com/greysquirr3l/lashes/internal/storage"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB(opts storage.Options) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch opts.Type {
	case storage.SQLite:
		dialector = sqlite.Open(opts.DSN)
	case storage.MySQL:
		dialector = mysql.Open(opts.DSN)
	case storage.Postgres:
		dialector = postgres.Open(opts.DSN)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", opts.Type)
	}

	config := &gorm.Config{}
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(opts.MaxConnections)
	sqlDB.SetConnMaxLifetime(opts.ConnTimeout)

	// Auto-migrate the schema
	if err := db.AutoMigrate(&ProxyModel{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
