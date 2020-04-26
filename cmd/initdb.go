package cmd

// initDB is a command wrapper for DB bootstrap
type initDB struct {
	common
}

// Execute implements flags.Commander
func (i *initDB) Execute(args []string) error {
	return migrate(i.DBType, i.ConnStr, true)
}
