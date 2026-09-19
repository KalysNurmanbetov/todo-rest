package domain

import "github.com/google/uuid"

type TaskId struct {
	value string
}

func NewTaskId() TaskId {
	return TaskId{value: uuid.NewString()}
}

func TaskIdFrom(v string) (TaskId, error) {
	if err := uuid.Validate(v); err != nil {
		return TaskId{}, ErrInvalidTaskId
	}
	return TaskId{value: v}, nil
}

func (t TaskId) Value() string {
	return t.value
}
