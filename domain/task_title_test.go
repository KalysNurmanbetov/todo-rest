package domain_test

import (
	"errors"
	"strings"
	"testing"
	"todo-rest/domain"
)

func TestNewTaskTitle_Normalization(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"valid title", "Buy some milk!", "Buy some milk!"},
		{"trimed value", "   Buy some milk!", "Buy some milk!"},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			res, err := domain.NewTaskTitle(it.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res.Value() != it.expected {
				t.Errorf("Value() = %q, want %q", res, it.expected)
			}
		})
	}
}

func TestNewTaskTitle_Invariants(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{"valid title", "Buy some milk", nil},
		{"empty  title", "", domain.ErrEmptyTitle},
		{"title with white spaces", "     ", domain.ErrEmptyTitle},
		{"valid title length", strings.Repeat("B", 200), nil},
		{"too long title", strings.Repeat("B", 201), domain.ErrTooLongTitle},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			_, err := domain.NewTaskTitle(it.input)
			if !errors.Is(err, it.expectedErr) {
				t.Errorf("NewTaskTitle(%q) error = %v, want %v", it.input, err, it.expectedErr)
			}
		})
	}
}
