package dto

import "github.com/jinzhu/gorm"

// Task is a scheduled task
type Task struct {
	gorm.Model

	// Description is the description of the task
	Description string
}
