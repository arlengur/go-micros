package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	errs "github.com/arlengur/go_assessment/internal/errors"
	"github.com/arlengur/go_assessment/internal/model"
	"github.com/arlengur/go_assessment/internal/repository/entity"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

var _ Repository = (*repository)(nil)

type repository struct {
	db *sql.DB
}

func NewRepository() *repository {
	db, err := sql.Open("sqlite", "./tasks.db")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return &repository{
		db: db,
	}
}

func (r *repository) CreateTask(ctx context.Context, taskType string, payload string) (string, error) {
	t := &model.Task{
		ID:        uuid.New().String(),
		Type:      taskType,
		Payload:   payload,
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
	}

	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO tasks (id, type, payload, status, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return "", errs.NewDbError("failed to prepare statement", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("failed to close statement: %v", err)
		}
	}()

	_, err = stmt.ExecContext(ctx, t.ID, t.Type, t.Payload, t.Status, t.CreatedAt)
	if err != nil {
		return "", errs.NewDbError("failed to execute statement", err)
	}

	return t.ID, nil
}

func (r *repository) GetTasks(ctx context.Context) ([]model.Task, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, type, payload, status, created_at, updated_at, error_info FROM tasks")
	if err != nil {
		return nil, errs.NewDbError("failed to query tasks", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("failed to close rows: %v", err)
		}
	}()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var t entity.TaskRecord
		if err := rows.Scan(&t.ID, &t.Type, &t.Payload, &t.Status, &t.CreatedAt, &t.UpdatedAt, &t.ErrorInfo); err != nil {
			return nil, err
		}
		task, err := entity.ToTask(t)
		if err != nil {
			return nil, fmt.Errorf("failed to convert task record to model: %v", err)
		}
		tasks = append(tasks, task)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("errors during iteration: %v", err)
	}
	return tasks, nil
}

func (r *repository) GetTaskByID(ctx context.Context, id string) (model.Task, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, type, payload, status, created_at, updated_at, error_info FROM tasks WHERE id = ?")
	if err != nil {
		return model.Task{}, errs.NewDbError("failed to prepare statement", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("failed to close statement: %v", err)
		}
	}()

	var t entity.TaskRecord
	err = stmt.QueryRowContext(ctx, id).Scan(&t.ID, &t.Type, &t.Payload, &t.Status, &t.CreatedAt, &t.UpdatedAt, &t.ErrorInfo)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, errs.ErrNotFound
	}
	if err != nil {
		return model.Task{}, errs.NewDbError("failed to execute statement", err)
	}

	task, err := entity.ToTask(t)
	if err != nil {
		return model.Task{}, fmt.Errorf("failed to convert task record to model: %v", err)
	}

	return task, nil
}

func (r *repository) DeleteTask(ctx context.Context, id string) error {
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM tasks WHERE id = ?")
	if err != nil {
		return errs.NewDbError("failed to prepare statement", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("failed to close statement: %v", err)
		}
	}()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return errs.NewDbError("failed to execute statement", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rows == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func (r *repository) GetPendingTasks(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id FROM tasks WHERE status = 'pending' LIMIT 20")
	if err != nil {
		return nil, errs.NewDbError("failed to query tasks", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *repository) UpdatePendingTask(ctx context.Context, id string) error {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE tasks SET status = 'processing', updated_at = ? WHERE id = ? AND status = 'pending'")
	if err != nil {
		return errs.NewDbError("failed to prepare statement", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("failed to close statement: %v", err)
		}
	}()

	res, err := stmt.ExecContext(ctx, time.Now(), id)
	if err != nil {
		return errs.NewDbError("failed to execute statement", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("task was cancelled or handled elsewhere")
	}

	return nil
}

func (r *repository) UpdateProcessingTask(ctx context.Context, status model.TaskStatus, errInfo string, id string) error {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE tasks SET status = ?, error_info = ?, updated_at = ? WHERE id = ? AND status = 'processing'")
	if err != nil {
		return errs.NewDbError("failed to prepare statement", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("failed to close statement: %v", err)
		}
	}()

	res, err := stmt.ExecContext(ctx, status, errInfo, time.Now(), id)
	if err != nil {
		return errs.NewDbError("failed to execute statement", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("task was cancelled or handled elsewhere")
	}

	return nil
}

func (r *repository) CancelTask(ctx context.Context, id string) error {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE tasks SET status = 'cancelled', updated_at = ? WHERE id = ? AND status IN ('pending', 'processing')")
	if err != nil {
		return errs.NewDbError("failed to prepare statement", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("failed to close statement: %v", err)
		}
	}()

	res, err := stmt.ExecContext(ctx, time.Now(), id)
	if err != nil {
		return errs.NewDbError("failed to execute statement", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("task was cancelled or handled elsewhere")
	}

	return nil
}
