package domain

import (
	"errors"
	"fmt"
)

// Invariants Errors
var ErrTaskAlreadyCompleted = errors.New("task is already completed")
var ErrTaskAlreadyUncompleted = errors.New("task is already uncompleted")
var ErrChangeTitleOfAlreadyCompleted = fmt.Errorf("title cannot be change %w", ErrTaskAlreadyCompleted)

// Fabric Errors
var ErrInvalidTaskId = errors.New("invalid task id")
var ErrEmptyTitle = errors.New("title must be not empty")
var ErrTooLongTitle = errors.New("title should contain max 200 characters")

// Repository Errors
var ErrTaskNotFound = errors.New("task is not found")
