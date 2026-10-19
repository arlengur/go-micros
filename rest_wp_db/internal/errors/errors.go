package errors

import (
	"errors"
)

var (
	ErrNotFound = errors.New("task not found")

	// Value Object errors
	ErrInvalidTaskStatus = errors.New("invalid task status")
)
