package model

import "time"

type Task struct {
	ID        string
	Type      string
	Payload   string
	Status    TaskStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	ErrorInfo string
}
