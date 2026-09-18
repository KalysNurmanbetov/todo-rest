package domain

import "strings"

type TaskTitle struct {
	value string
}

func NewTaskTitle(v string) (TaskTitle, error) {
	if strings.Trim(v, " ") == "" {
		return TaskTitle{}, ErrEmptyTitle
	}
	if len(v) > 200 {
		return TaskTitle{}, ErrTooLongTitle
	}

	return TaskTitle{value: v}, nil
}

func (t TaskTitle) Value() string {
	return t.value
}
