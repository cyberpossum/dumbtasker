package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/cyberpossum/dumbtasker/dal"
	"github.com/cyberpossum/dumbtasker/dto"
	"github.com/cyberpossum/tabwriter"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
)

const (
	_redCutoff    = 2 * time.Hour
	_yellowCutoff = 6 * time.Hour
)

type listTask struct {
	common

	Due      smartTime `long:"due" description:"Tasks with due date before the specified"`
	All      bool      `short:"a" long:"all" description:"Show all tasks including completed"`
	NoColor  bool      `long:"no-color" description:"suppress color output"`
	FullDate bool      `long:"full-date" description:"Show full date in output"`
}

func (l *listTask) formatOutput(t *dto.Task) (string, string) {
	idStr := fmt.Sprintf("%08d", t.ID)
	var dueStr string
	if l.FullDate {
		dueStr = fmt.Sprintf("%v", t.Due)
	} else {
		dueStr = humanize.Time(t.Due)
	}

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
	due := time.Time(l.Due)

	tasks, err := dal.ListTasks(l.getDbConfig(), due, l.All)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 4, 1, byte(' '), tabwriter.ANSIColors)
	for _, tt := range tasks {
		id, fdue := l.formatOutput(&tt)
		fmt.Fprintf(w, "%v\t%q:\t%v\n", id, tt.Description, fdue)
	}
	w.Flush()

	return nil
}
