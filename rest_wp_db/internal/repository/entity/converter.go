package entity

import "github.com/arlengur/go_assessment/internal/model"

func ToTask(r TaskRecord) (model.Task, error) {
	taskStatus, err := model.NewTaskStatus(r.Status)
	if err != nil {
		return model.Task{}, err
	}

	return model.Task{
		ID:        r.ID,
		Type:      r.Type,
		Payload:   r.Payload,
		Status:    taskStatus,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		ErrorInfo: r.ErrorInfo.String,
	}, nil
}
