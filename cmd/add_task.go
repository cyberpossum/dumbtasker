package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/cyberpossum/dumbtasker/database"
	"github.com/cyberpossum/dumbtasker/dto"
	"github.com/jinzhu/gorm"
)

type smartTime time.Time

func (st *smartTime) UnmarshalFlag(value string) error {
	t, err := dateparse.ParseLocal(value)
	if err != nil {
		return err
	}

	*st = smartTime(t)

	return nil
}

type addTask struct {
	common

	Due      smartTime     `long:"due" description:"Task due date" required:"true"`
	Estimate time.Duration `long:"estimate" description:"task duration estimate"`

	PosArgs struct {
		Description string `positional-arg-name:"description" description:"task description, e.g. \"buy food\""`
	} `positional-args:"yes" required:"yes"`
}

func (a *addTask) newTask() (*dto.Task, error) {
	trimmedDesc := strings.Trim(a.PosArgs.Description, " \t\n")
	if len(trimmedDesc) == 0 {
		return nil, errors.New("empty task description")
	}

	if len(trimmedDesc) > 255 {
		return nil, errors.New("description too long, max length is 255 characters")
	}

	dueTime := time.Time(a.Due)

	if dueTime.Before(time.Now()) {
		return nil, errors.New("cannot schedule task in the past")
	}

	return &dto.Task{
		Model:       gorm.Model{},
		Description: trimmedDesc,
		Due:         dueTime,
		Estimate:    a.Estimate,
		Status:      dto.Open,
	}, nil
}

// Execute executes the command, adding a task
func (a *addTask) Execute([]string) error {
	task, err := a.newTask()
	if err != nil {
		return fmt.Errorf("creating task: %w", err)
	}

	db, err := database.OpenDB(a.DBType, a.ConnStr)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	return db.Create(task).Error
}
