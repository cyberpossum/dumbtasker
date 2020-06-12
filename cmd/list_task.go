package cmd

import (
	"fmt"
	"time"

	"github.com/cyberpossum/dumbtasker/database"
	"github.com/cyberpossum/dumbtasker/dto"
)

type listTask struct {
	common

	Due smartTime `long:"due" description:"Tasks with due date before the specified"`
	All bool      `short:"a" long:"all" description:"Show all tasks including completed"`
}

// Execute implements flags.Commander
func (l *listTask) Execute(args []string) error {
	db, err := database.OpenDB(l.DBType, l.ConnStr)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	var tasks []dto.Task

	due := time.Time(l.Due)

	db = db.Where(&dto.Task{Status: dto.Open})

	if l.All {
		db = db.Or(&dto.Task{Status: dto.Closed})
	}

	if !due.IsZero() {
		db = db.Where("due <= ?", due)
	}

	if err := db.Find(&tasks).Error; err != nil {
		return err
	}

	for _, tt := range tasks {
		fmt.Printf("%08d\t%q:\t%v\n", tt.ID, tt.Description, tt.Due)
	}

	return nil
}
