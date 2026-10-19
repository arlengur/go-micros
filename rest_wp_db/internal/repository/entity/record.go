package entity

import (
	"database/sql"
	"time"
)

type TaskRecord struct {
	ID        string
	Type      string
	Payload   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	ErrorInfo sql.NullString
}
