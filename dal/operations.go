package dal

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cyberpossum/dumbtasker/database"
	"github.com/cyberpossum/dumbtasker/dto"
)

// CreateTask creates the task in the database
func CreateTask(cfg *database.Config, t *dto.Task) error {
	db, err := database.OpenDB(cfg)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	return db.Create(t).Error
}

func isExpectedStatus(status dto.TaskStatus, expectedStatuses []dto.TaskStatus) bool {
	for _, ts := range expectedStatuses {
		if ts == status {
			return true
		}
	}

	return false
}

// ChangeTaskStatus changes the status of the task in the database
func ChangeTaskStatus(cfg *database.Config, taskID string, expectedStatuses []dto.TaskStatus, newStatus dto.TaskStatus) error {
	id, err := strconv.Atoi(taskID)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	db, err := database.OpenDB(cfg)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	var t dto.Task

	if err := db.First(&t, id).Error; err != nil {
		return fmt.Errorf("error looking up record with ID %d: %w", id, err)
	}

	if !isExpectedStatus(t.Status, expectedStatuses) {
		return fmt.Errorf("unexpected status for task with ID %d: %v", id, t.Status)
	}

	t.Status = newStatus
	if err := db.Save(t).Error; err != nil {
		return fmt.Errorf("updating record: %w", err)
	}

	return nil
}

// Migrate perform database schema migration to the current version
// if `checkIfEmpty` argument is set, it will fail on a non-empty database
func Migrate(cfg *database.Config, checkIfEmpty bool) error {
	db, err := database.OpenDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if checkIfEmpty {
		if err := database.EnsureEmpty(db, cfg.Type); err != nil {
			return err
		}
	}

	db.AutoMigrate(dto.AllModels...)

	return nil
}

// ListTasks lists the task matching search criteria
// Tasks are ordered by due date, ascending
func ListTasks(cfg *database.Config, due time.Time, allTasks bool) ([]dto.Task, error) {
	db, err := database.OpenDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	var tasks []dto.Task

	db = db.Where(&dto.Task{Status: dto.Open})

	if allTasks {
		db = db.Or(&dto.Task{Status: dto.Closed})
	}

	if !due.IsZero() {
		db = db.Where("due <= ?", due)
	}

	db = db.Order("due ASC")

	if err := db.Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("listing tasks: %w", err)
	}

	return tasks, nil
}
