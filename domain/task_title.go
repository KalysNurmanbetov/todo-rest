package domain

import "strings"

type TaskTitle struct {
	value string
}

func NewTaskTitle(v string) (TaskTitle, error) {
	if trimedV := strings.Trim(v, " "); trimedV == "" {
		return TaskTitle{}, ErrEmptyTitle
	} else {
		if len(trimedV) > 200 {
			return TaskTitle{}, ErrTooLongTitle
		}
		return TaskTitle{value: trimedV}, nil
	}
}

func (t TaskTitle) Value() string {
	return t.value
}
