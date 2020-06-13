package cmd

import (
	"github.com/cyberpossum/dumbtasker/dal"
	"github.com/cyberpossum/dumbtasker/dto"
)

type deleteTask struct {
	common

	PosArgs struct {
		ID string `positional-arg-name:"id" description:"task ID, e.g. \"12\" or \"00000012\""`
	} `positional-args:"yes" required:"yes"`
}

// Execute executes the command, adding a task
func (d *deleteTask) Execute([]string) error {
	return dal.ChangeTaskStatus(d.getDbConfig(), d.PosArgs.ID, []dto.TaskStatus{dto.Open, dto.Closed}, dto.Deleted)
}
