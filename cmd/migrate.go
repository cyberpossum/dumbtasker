package cmd

import (
	"github.com/cyberpossum/dumbtasker/database"
	"github.com/cyberpossum/dumbtasker/dto"
)

// migrateDB is a command wrapper for DB migration
type migrateDB struct {
	common
}

func migrate(dbType database.Type, connStr string, checkIfEmpty bool) error {
	db, err := database.OpenDB(dbType, connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	if checkIfEmpty {
		if err := database.EnsureEmpty(db, dbType); err != nil {
			return err
		}
	}

	db.AutoMigrate(dto.AllModels...)

	return nil
}

// Execute implements flags.Commander
func (i *migrateDB) Execute(args []string) error {
	return migrate(i.DBType, i.ConnStr, false)
}
