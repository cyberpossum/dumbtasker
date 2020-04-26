package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

func selectStrings(db *gorm.DB, selectStmt string) ([]string, error) {
	var strValues []string
	if err := db.Raw(selectStmt).Pluck("table_nane", &strValues).Error; err != nil {
		return nil, err
	}
	return strValues, nil
}

// EnsureEmpty checks if the database has no tables
func EnsureEmpty(db *gorm.DB, dbType Type) error {
	selectStmt, ok := tableSelectStmts[dbType]
	if !ok {
		return fmt.Errorf("error selecting table list: unknown database type: %v", dbType)
	}

	tables, err := selectStrings(db, selectStmt)
	if err != nil {
		return fmt.Errorf("error selecting table list: %w", err)
	}

	if len(tables) > 0 {
		return fmt.Errorf("database not empty, at least one table present: %q", tables[0])
	}

	return nil
}
