package cmd

import (
	"github.com/cyberpossum/dumbtasker/database"
	"github.com/jessevdk/go-flags"
)

type common struct {
	ConnStr string        `long:"connection-string" description:"Connection string" default:"file:dumbtasker.db?cache=shared" env:"DUMBTASKER_CONNSTR"`
	DBType  database.Type `long:"db-type" description:"Database type" default:"sqlite3" env:"DUMBTASKER_DBTYPE"`
}

// Opts is the list of commands and options
type Opts struct {
	InitDB    *initDB    `command:"init-db" description:"Initialize an empty database"`
	MigrateDB *migrateDB `command:"migrate-db" description:"Migrate existing database"`
	AddTask   *addTask   `command:"add" description:"Add a new task" alias:"new" alias:"new-task"`
}

// AppendDBTypes populates choices for DBType option
func AppendDBTypes(p *flags.Parser) {
	for _, c := range p.Commands() {
		if dbType := c.FindOptionByLongName("db-type"); dbType != nil {
			dbType.Choices = database.PossibleDBValues()
		}
	}
}
