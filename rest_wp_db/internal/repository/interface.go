package repository

import (
	"context"

	"github.com/arlengur/go_assessment/internal/model"
)

type Repository interface {
	CreateTask(ctx context.Context, taskType string, payload string) (string, error)
	GetTasks(ctx context.Context) ([]model.Task, error)
	GetTaskByID(ctx context.Context, id string) (model.Task, error)
	CancelTask(ctx context.Context, id string) error
	DeleteTask(ctx context.Context, id string) error
	GetPendingTasks(ctx context.Context) ([]string, error)
	UpdatePendingTask(ctx context.Context, id string) error
	UpdateProcessingTask(ctx context.Context, status model.TaskStatus, err string, id string) error
}
