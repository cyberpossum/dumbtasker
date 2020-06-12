package cmd

import (
	"fmt"
	"strconv"

	"github.com/cyberpossum/dumbtasker/database"
	"github.com/cyberpossum/dumbtasker/dto"
)

func isExpectedStatus(status dto.TaskStatus, expectedStatuses []dto.TaskStatus) bool {
	for _, ts := range expectedStatuses {
		if ts == status {
			return true
		}
	}

	return false
}

func changeTaskStatus(dbType database.Type, connString string, taskID string, expectedStatuses []dto.TaskStatus, newStatus dto.TaskStatus) error {
	id, err := strconv.Atoi(taskID)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	db, err := database.OpenDB(dbType, connString)
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
