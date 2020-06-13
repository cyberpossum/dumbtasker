package cmd

import "github.com/cyberpossum/dumbtasker/dal"

// migrateDB is a command wrapper for DB migration
type migrateDB struct {
	common
}

// Execute implements flags.Commander
func (i *migrateDB) Execute(args []string) error {
	return dal.Migrate(i.getDbConfig(), false)
}
