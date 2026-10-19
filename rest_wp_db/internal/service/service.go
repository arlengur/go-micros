package service

import (
	"context"

	"github.com/arlengur/go_assessment/internal/model"
	"github.com/arlengur/go_assessment/internal/pool"
	"github.com/arlengur/go_assessment/internal/repository"
)

var _ Service = (*taskService)(nil)

type taskService struct {
	pool pool.Pool
	repo repository.Repository
}

func NewTaskService(pool pool.Pool, repo repository.Repository) *taskService {
	ts := &taskService{
		pool: pool,
		repo: repo,
	}
	return ts
}

func (s *taskService) CreateTask(ctx context.Context, taskType string, payload string) (string, error) {
	uuid, err := s.repo.CreateTask(ctx, taskType, payload)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (s *taskService) GetTasks(ctx context.Context) ([]model.Task, error) {
	return s.repo.GetTasks(ctx)
}

func (s *taskService) GetTaskByID(ctx context.Context, id string) (model.Task, error) {
	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (s *taskService) CancelTask(ctx context.Context, id string) error {
	return s.pool.CancelTask(id)
}

func (s *taskService) DeleteTask(ctx context.Context, id string) error {
	return s.repo.DeleteTask(ctx, id)
}
