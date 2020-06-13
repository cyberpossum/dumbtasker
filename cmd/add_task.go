package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/cyberpossum/dumbtasker/dal"
	"github.com/cyberpossum/dumbtasker/dto"
	"github.com/jinzhu/gorm"
)

type smartTime time.Time

func (st *smartTime) UnmarshalFlag(value string) error {
	// formatting
	// "today", "tomorrow" -> due until EOD
	// "today 14:00" -> until today's 14:00
	procValue := strings.TrimSpace(value)

	switch procValue {
	case "today":
		procValue = time.Now().Add(24 * time.Hour).Format("2006-01-02")
	case "tomorrow":
		procValue = time.Now().Add(48 * time.Hour).Format("2006-01-02")
	}
	procValue = strings.Replace(procValue, "today", time.Now().Format("2006-01-02"), 1)
	procValue = strings.Replace(procValue, "tomorrow", time.Now().Add(24*time.Hour).Format("2006-01-02"), 1)

	t, err := dateparse.ParseLocal(procValue)
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

	return dal.CreateTask(a.getDbConfig(), task)
}
