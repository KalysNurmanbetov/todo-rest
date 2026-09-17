package domain

import (
	"strings"
	"time"
)

type Task struct {
	Id          string
	Title       string
	Completed   bool
	CreateAt    time.Time
	CompletedAt *time.Time
}

func NewTask(id string, title string) (Task, error) {
	if strings.Trim(title, " ") == "" {
		return Task{}, ErrEmptyTitle
	}

	return Task{
		Id:          id,
		Title:       title,
		Completed:   false,
		CreateAt:    time.Now(),
		CompletedAt: nil,
	}, nil
}

func (t *Task) ChangeTitle(title string) error {
	t.Title = title
	return nil
}

func (t *Task) Complete() error {
	if t.Completed {
		return ErrTaskAlreadyCompleted
	}
	completeTime := time.Now()

	t.Completed = true
	t.CompletedAt = &completeTime
	return nil
}

func (t *Task) Uncomplete() error {
	if t.Completed == false {
		return ErrTaskAlreadyUncompleted
	}
	t.Completed = false
	t.CompletedAt = nil
	return nil
}
