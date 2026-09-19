package domain_test

import (
	"errors"
	"testing"
	"time"
	"todo-rest/domain"
)

func TestNewTask(t *testing.T) {
	before := time.Now()
	task := newTestTask(t)
	after := time.Now()

	if task.CreatedAt().Before(before) || task.CreatedAt().After(after) {
		t.Errorf("CreatedAt() = %v, want between %v and %v", task.CreatedAt(), before, after)
	}
}

func TestTask_ChangeTitle_BeforeComplete(t *testing.T) {
	task := newTestTask(t)
	title := newTaskTitle(t, "Finish homework")

	err := task.ChangeTitle(title)
	if err != nil {
		t.Fatalf("ChangeTitle() unexpected error: %v", err)
	}

	if got := task.Title().Value(); got != title.Value() {
		t.Errorf("Title().Value() = %q, want %q", got, title.Value())
	}
}

func TestTask_ChangeTitle_AfterCompletion(t *testing.T) {

	task := newTestTask(t)
	task.Complete()
	title := newTaskTitle(t, "Finish homework")

	err := task.ChangeTitle(title)
	if !errors.Is(err, domain.ErrChangeTitleOfAlreadyCompleted) {
		t.Fatalf("ChangeTitle() must return: %v if task already completed", err)
	}

	if task.Title().Value() == title.Value() {
		t.Errorf("Title must not be changed after completion")
	}
}

func TestTask_Complete(t *testing.T) {
	task := newTestTask(t)

	before := time.Now()
	err := task.Complete()
	after := time.Now()

	if err != nil {
		t.Fatalf("Complete() unexpected error: %v", err)
	}
	if !task.Completed() {
		t.Errorf("Completed() = false, want true")
	}

	completedAt := task.CompletedAt()
	if completedAt == nil {
		t.Fatalf("CompletedAt() = nil, want non-nil")
	}
	if completedAt.Before(before) || completedAt.After(after) {
		t.Errorf("CompletedAt() = %v, want between %v and %v", completedAt, before, after)
	}
}

func TestTask_Complete_AlreadyCompleted(t *testing.T) {
	task := newTestTask(t)
	if err := task.Complete(); err != nil {
		t.Fatalf("Complete() unexpected error: %v", err)
	}
	completedAt := task.CompletedAt()

	err := task.Complete()
	if !errors.Is(err, domain.ErrTaskAlreadyCompleted) {
		t.Fatalf("Complete() error = %v, want %v", err, domain.ErrTaskAlreadyCompleted)
	}
	if task.CompletedAt() == nil || !task.CompletedAt().Equal(*completedAt) {
		t.Errorf("CompletedAt() changed on second Complete() call: got %v, want unchanged %v", task.CompletedAt(), completedAt)
	}
}

func TestTask_Uncomplete_BeforeComplete(t *testing.T) {
	task := newTestTask(t)

	err := task.Uncomplete()
	if !errors.Is(err, domain.ErrTaskAlreadyUncompleted) {
		t.Fatalf("Uncomplete() error = %v, want %v", err, domain.ErrTaskAlreadyUncompleted)
	}
}

func TestTask_Uncomplete_AfterComplete(t *testing.T) {
	task := newTestTask(t)
	if err := task.Complete(); err != nil {
		t.Fatalf("Complete() unexpected error: %v", err)
	}

	err := task.Uncomplete()
	if err != nil {
		t.Fatalf("Uncomplete() unexpected error: %v", err)
	}
	if task.Completed() {
		t.Errorf("Completed() = true, want false")
	}
	if task.CompletedAt() != nil {
		t.Errorf("CompletedAt() = %v, want nil", task.CompletedAt())
	}
}

func TestTask_Uncomplete_AlreadyUncompleted(t *testing.T) {
	task := newTestTask(t)
	if err := task.Complete(); err != nil {
		t.Fatalf("Complete() unexpected error: %v", err)
	}
	if err := task.Uncomplete(); err != nil {
		t.Fatalf("Uncomplete() unexpected error: %v", err)
	}

	err := task.Uncomplete()
	if !errors.Is(err, domain.ErrTaskAlreadyUncompleted) {
		t.Fatalf("Uncomplete() error = %v, want %v", err, domain.ErrTaskAlreadyUncompleted)
	}
}

func TestRehydrateTask(t *testing.T) {
	id := domain.NewTaskId()
	title := newTaskTitle(t, "Buy some milk!")
	createdAt := time.Now().Add(-24 * time.Hour)
	completedAt := time.Now().Add(-1 * time.Hour)

	tests := []struct {
		name        string
		completed   bool
		completedAt *time.Time
	}{
		{"not completed", false, nil},
		{"completed", true, &completedAt},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			task := domain.RehydrateTask(id, title, it.completed, createdAt, it.completedAt)

			if task.Id() != id {
				t.Errorf("Id() = %v, want %v", task.Id(), id)
			}
			if task.Title().Value() != title.Value() {
				t.Errorf("Title().Value() = %q, want %q", task.Title().Value(), title.Value())
			}
			if task.Completed() != it.completed {
				t.Errorf("Completed() = %v, want %v", task.Completed(), it.completed)
			}
			if !task.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", task.CreatedAt(), createdAt)
			}
			if (task.CompletedAt() == nil) != (it.completedAt == nil) {
				t.Errorf("CompletedAt() = %v, want %v", task.CompletedAt(), it.completedAt)
			} else if it.completedAt != nil && !task.CompletedAt().Equal(*it.completedAt) {
				t.Errorf("CompletedAt() = %v, want %v", *task.CompletedAt(), *it.completedAt)
			}
		})
	}
}

//Helpers

func newTestTask(t *testing.T) *domain.Task {
	t.Helper()
	title := newTaskTitle(t, "Buy some milk!")
	task := domain.NewTask(domain.NewTaskId(), title)
	return task
}

func newTaskTitle(t *testing.T, v string) domain.TaskTitle {
	t.Helper()
	title, err := domain.NewTaskTitle(v)
	if err != nil {
		t.Fatalf("NewTaskTitle: unexpected error: %v", err)
	}
	return title
}
