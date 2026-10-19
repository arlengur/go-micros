package model

import errs "github.com/arlengur/go_assessment/internal/errors"

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
	StatusCancelled  TaskStatus = "cancelled"
)

func NewTaskStatus(status string) (TaskStatus, error) {
	taskStatus := TaskStatus(status)
	switch taskStatus {
	case StatusPending, StatusProcessing, StatusCompleted, StatusFailed, StatusCancelled:
		return taskStatus, nil
	default:
		return "", errs.ErrInvalidTaskStatus
	}
}
