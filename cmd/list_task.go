package cmd

import (
	"fmt"
	"time"

	"github.com/cyberpossum/dumbtasker/database"
	"github.com/cyberpossum/dumbtasker/dto"
	"github.com/fatih/color"
)

const (
	_redCutoff    = 2 * time.Hour
	_yellowCutoff = 6 * time.Hour
)

type listTask struct {
	common

	Due     smartTime `long:"due" description:"Tasks with due date before the specified"`
	All     bool      `short:"a" long:"all" description:"Show all tasks including completed"`
	NoColor bool      `long:"no-color" description:"suppress color output"`
}

func (l *listTask) formatOutput(t *dto.Task) (string, string) {
	idStr := fmt.Sprintf("%08d", t.ID)
	dueStr := fmt.Sprintf("%v", t.Due)
	if l.NoColor {
		return idStr, dueStr
	}
	dueFmt := color.GreenString
	switch {
	case t.Due.Before(time.Now()):
		dueFmt = color.HiRedString
	case t.Due.Before(time.Now().Add(_redCutoff)):
		dueFmt = color.RedString
	case t.Due.Before(time.Now().Add(_yellowCutoff)):
		dueFmt = color.YellowString
	}

	if t.Status == dto.Closed {
		idStr = color.HiBlackString(idStr)
	}

	return idStr, dueFmt(dueStr)
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
		id, fdue := l.formatOutput(&tt)
		fmt.Printf("%v\t%q:\t%v\n", id, tt.Description, fdue)
	}

	return nil
}
