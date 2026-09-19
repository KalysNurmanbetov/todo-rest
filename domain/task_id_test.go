package domain_test

import (
	"errors"
	"testing"
	"todo-rest/domain"

	"github.com/google/uuid"
)

func TestNewTaskId(t *testing.T) {
	id := domain.NewTaskId()

	if id.Value() == "" {
		t.Fatal("NewTaskId() returned an empty id")
	}

	if err := uuid.Validate(id.Value()); err != nil {
		t.Errorf("NewTaskId() = %q, not a valid uuid: %v", id.Value(), err)
	}
}

func TestNewTaskId_Unique(t *testing.T) {
	first := domain.NewTaskId()
	second := domain.NewTaskId()

	if first.Value() == second.Value() {
		t.Errorf("NewTaskId() called twice returned the same value: %q", first.Value())
	}
}

func TestTaskIdFrom(t *testing.T) {
	const validId string = "bef65b1c-e4f6-451d-a8cd-f5fadd3db037"

	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{"valid id", validId, nil},
		{"invalid id", "some_id", domain.ErrInvalidTaskId},
		{"empty id", "", domain.ErrInvalidTaskId},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			_, err := domain.TaskIdFrom(it.input)
			if !errors.Is(err, it.expectedErr) {
				t.Errorf("TaskIdFrom(%q) error = %v, want %v", it.input, err, it.expectedErr)
			}
		})
	}
}
