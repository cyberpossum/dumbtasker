package cmd

import (
	"github.com/cyberpossum/dumbtasker/database"
	"github.com/jessevdk/go-flags"
)

type common struct {
	ConnStr string        `long:"connection-string" description:"Connection string" default:"file:dumbtasker.db?cache=shared" env:"DUMBTASKER_CONNSTR"`
	DBType  database.Type `long:"db-type" description:"Database type" default:"sqlite3" env:"DUMBTASKER_DBTYPE"`
}

func (c *common) getDbConfig() *database.Config {
	return &database.Config{
		Type:    c.DBType,
		ConnStr: c.ConnStr,
	}
}

// Opts is the list of commands and options
type Opts struct {
	InitDB       *initDB       `command:"init-db" description:"Initialize an empty database"`
	MigrateDB    *migrateDB    `command:"migrate-db" description:"Migrate existing database"`
	AddTask      *addTask      `command:"add" description:"Add a new task" alias:"new" alias:"new-task"`
	ListTask     *listTask     `command:"list" description:"List tasks"`
	CompleteTask *completeTask `command:"done" description:"Complete a task" alias:"complete" alias:"close"`
	DeleteTask   *deleteTask   `command:"delete" description:"Delete a task" alias:"cancel" alias:"remove" alias:"erase"`
}

// AppendDBTypes populates choices for DBType option
func AppendDBTypes(p *flags.Parser) {
	for _, c := range p.Commands() {
		if dbType := c.FindOptionByLongName("db-type"); dbType != nil {
			dbType.Choices = database.PossibleDBValues()
		}
	}
}
