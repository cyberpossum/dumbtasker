package cmd

import (
	"fmt"
	"strconv"

	"github.com/cyberpossum/dumbtasker/database"
	"github.com/cyberpossum/dumbtasker/dto"
)

type completeTask struct {
	common

	PosArgs struct {
		ID string `positional-arg-name:"id" description:"task ID, e.g. \"12\" or \"00000012\""`
	} `positional-args:"yes" required:"yes"`
}

// Execute executes the command, adding a task
func (a *completeTask) Execute([]string) error {
	id, err := strconv.Atoi(a.PosArgs.ID)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	db, err := database.OpenDB(a.DBType, a.ConnStr)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	var t dto.Task

	if err := db.First(&t, id).Error; err != nil {
		return fmt.Errorf("error looking up record with ID %d: %w", id, err)
	}

	t.Status = dto.Closed
	if err := db.Save(t).Error; err != nil {
		return fmt.Errorf("updating record: %w", err)
	}

	return nil
}
