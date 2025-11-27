package helper

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrDuplicateEntry    = errors.New("duplicate entry")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternalServer    = errors.New("internal server error")
	ErrDatabaseOperation = errors.New("database operation failed")
)

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrNotFound)
}

// WrapDatabaseError wraps database errors with custom error types
func WrapDatabaseError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return ErrDatabaseOperation
}
