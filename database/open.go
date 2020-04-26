package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// OpenDB opens a gorm database connection
func OpenDB(dbType Type, connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(dbType.String(), connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	return db, nil
}
