package dto

import (
	"time"

	"github.com/jinzhu/gorm"
)

// TaskStatus is the status for the task
// ses const values for details
type TaskStatus uint64

const (
	// Open is the status for active, planned tasks
	Open TaskStatus = iota + 1
	// Closed is the status for tasks that had been completed
	Closed
	// Deleted is the status for deleted tasks
	Deleted
)

// Task is a scheduled task
type Task struct {
	gorm.Model

	// Description is the description of the task
	Description string `gorm:"size:255"`
	// Due is the due date for the task
	Due time.Time `gorm:"index:ind_due"`
	// Estimate is the estimated time to complete the task
	Estimate time.Duration
	// Status
	Status TaskStatus `gorm:"index:ind_status"`
}
