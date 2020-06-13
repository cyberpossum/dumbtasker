package cmd

import "github.com/cyberpossum/dumbtasker/dal"

// initDB is a command wrapper for DB bootstrap
type initDB struct {
	common
}

// Execute implements flags.Commander
func (i *initDB) Execute(args []string) error {
	return dal.Migrate(i.getDbConfig(), true)
}
