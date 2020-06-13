package database

import (
	"fmt"
	"strings"

	"github.com/jessevdk/go-flags"
)

// Type is the database used to store the data
type Type uint16

const (
	// DBTypeSQLite3 represent SQLite3 database
	DBTypeSQLite3 Type = iota + 1
	// DBTypeMySQL represent MySQL database
	DBTypeMySQL
)

// Config contains all the necessary information to
// access the database
type Config struct {
	Type
	ConnStr string
}

var dbTypeMapping map[string]Type = map[string]Type{
	"sqlite3": DBTypeSQLite3,
	"mysql":   DBTypeMySQL,
}

var tableSelectStmts = map[Type]string{
	DBTypeMySQL:   "SELECT `table_name` FROM `information_schema`.`tables` WHERE `TABLE_SCHEMA` <> 'information_schema';",
	DBTypeSQLite3: `SELECT "name" AS "table_name" FROM "sqlite_master" WHERE "type" ='table' AND "name" NOT LIKE 'sqlite_%';`,
}

// UnmarshalFlag implements flags.Unmarshaler interface
func (d *Type) UnmarshalFlag(value string) error {
	dbVal, ok := dbTypeMapping[value]
	if !ok {
		return fmt.Errorf("error unmarshalling DatabaseType: invalid value %q", value)
	}
	*d = dbVal
	return nil
}

// String implements fmt.Stringer
func (d Type) String() string {
	text, err := d.MarshalFlag()
	if err != nil {
		return "ERR: " + err.Error()
	}
	return text
}

// MarshalFlag implements flags.Marshaler interface
func (d Type) MarshalFlag() (string, error) {
	for typeStr, typeVal := range dbTypeMapping {
		if typeVal == d {
			return typeStr, nil
		}
	}
	return "", fmt.Errorf("error unmarshalling DatabaseType: invalid value %v", d)
}

// Complete implements flags.Completer interface
func (d *Type) Complete(match string) []flags.Completion {
	var result []flags.Completion

	for _, typeStr := range PossibleDBValues() {
		if strings.HasPrefix(typeStr, match) {
			result = append(result, flags.Completion{Item: typeStr})
		}
	}

	return result
}

// PossibleDBValues returns the list of possible values for DB key
func PossibleDBValues() []string {
	var result []string

	for typeStr := range dbTypeMapping {
		result = append(result, typeStr)
	}

	return result
}
