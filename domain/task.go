package domain

import (
	"time"
)

type Task struct {
	id          TaskId
	title       TaskTitle
	completed   bool
	createAt    time.Time
	completedAt *time.Time
}

func NewTask(id TaskId, title TaskTitle) *Task {
	return &Task{
		id:          id,
		title:       title,
		completed:   false,
		createAt:    time.Now(),
		completedAt: nil,
	}
}

func (t *Task) ChangeTitle(title TaskTitle) error {
	if t.completed {
		return ErrChangeTitleOfAlreadyCompleted
	}
	t.title = title
	return nil
}

func (t *Task) Complete() error {
	if t.completed {
		return ErrTaskAlreadyCompleted
	}
	completeTime := time.Now()

	t.completed = true
	t.completedAt = &completeTime
	return nil
}

func (t *Task) Uncomplete() error {
	if t.completed == false {
		return ErrTaskAlreadyUncompleted
	}
	t.completed = false
	t.completedAt = nil
	return nil
}

func (t *Task) Id() TaskId {
	return t.id
}

func (t *Task) Title() TaskTitle {
	return t.title
}

func (t *Task) Completed() bool {
	return t.completed
}

func (t *Task) CreatedAt() time.Time {
	return t.createAt
}

func (t *Task) CompletedAt() *time.Time {
	return t.completedAt
}

// Only for infrustructure level to restore the state from storage
func RehydrateTask(id TaskId, title TaskTitle, completed bool, createdAt time.Time, completedAt *time.Time) *Task {
	return &Task{
		id: id, title: title, completed: completed, createAt: createdAt, completedAt: completedAt,
	}
}
