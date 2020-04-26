package main

import (
	"fmt"
	"os"

	"github.com/cyberpossum/dumbtasker/cmd"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jessevdk/go-flags"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func exitError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func main() {
	var options cmd.Opts

	var parser = flags.NewParser(&options, flags.Default)

	// Append DB type choices
	cmd.AppendDBTypes(parser)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
