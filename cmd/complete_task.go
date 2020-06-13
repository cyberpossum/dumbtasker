package cmd

import (
	"github.com/cyberpossum/dumbtasker/dal"
	"github.com/cyberpossum/dumbtasker/dto"
)

type completeTask struct {
	common

	PosArgs struct {
		ID string `positional-arg-name:"id" description:"task ID, e.g. \"12\" or \"00000012\""`
	} `positional-args:"yes" required:"yes"`
}

// Execute executes the command, adding a task
func (c *completeTask) Execute([]string) error {
	return dal.ChangeTaskStatus(c.getDbConfig(), c.PosArgs.ID, []dto.TaskStatus{dto.Open}, dto.Closed)
}
