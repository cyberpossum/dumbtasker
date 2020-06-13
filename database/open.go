package database

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

// OpenDB opens a gorm database connection
func OpenDB(cfg *Config) (*gorm.DB, error) {
	if cfg == nil {
		return nil, errors.New("nil config provided for OpenDB")
	}
	db, err := gorm.Open(cfg.Type.String(), cfg.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	return db, nil
}
