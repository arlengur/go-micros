package service

import (
	"context"

	"github.com/arlengur/go_assessment/internal/model"
)

type Service interface {
	CreateTask(ctx context.Context, taskType string, payload string) (string, error)
	GetTasks(ctx context.Context) ([]model.Task, error)
	GetTaskByID(ctx context.Context, id string) (model.Task, error)
	CancelTask(ctx context.Context, id string) error
	DeleteTask(ctx context.Context, id string) error
}
