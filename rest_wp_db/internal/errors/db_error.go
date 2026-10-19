package errors

import "fmt"

type DbError struct {
	Message string
	Err     error
}

func NewDbError(message string, err error) *DbError {
	return &DbError{
		Message: message,
		Err:     err,
	}
}

func (e *DbError) Error() string {
	return fmt.Sprintf("%v: %v", e.Message, e.Err)
}
