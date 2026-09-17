package domain

import (
	"errors"
)

// Invariants Errors
var ErrTaskAlreadyCompleted = errors.New("task is already completed")
var ErrTaskAlreadyUncompleted = errors.New("task is already uncompleted")

// Fabric Errors
var ErrEmptyTitle = errors.New("title must be not empty")

// Repository Errors
var ErrTaskNotFound = errors.New("task is not found")
