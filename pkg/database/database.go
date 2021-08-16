package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New creates a new database connection
func New(dialect string, dsn string, config *gorm.Config) (*gorm.DB, error) {
	if config == nil {
		config = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}

	if dialect == "sqlite" {
		config.DisableForeignKeyConstraintWhenMigrating = true
		return NewSQLite(dsn, config)
	} else if dialect == "mysql" {
		return NewMySQL(dsn, config)
	}

	return nil, fmt.Errorf("unsupported dialect: %s", dialect)
}

// NewMySQL Creates a new MySQL Database Connection
func NewMySQL(dsn string, config *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), config)
}

// NewSQLite --
func NewSQLite(file string, config *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(file), config)
}
