package domain

import "fmt"

type TaskAlreadCompletedError struct {
	TaskId string
}

func (e TaskAlreadCompletedError) Error() string {
	return fmt.Sprintf("task with id '%s' is already compeleted", e.TaskId)
}

type TaskAlreadUncompletedError struct {
	TaskId string
}

func (e TaskAlreadUncompletedError) Error() string {
	return fmt.Sprintf("task with id '%s' is already uncompeleted", e.TaskId)

}

type TaskNotFoundError struct {
	TaskId string
}

func (e TaskNotFoundError) Error() string {
	return fmt.Sprintf("task with id '%s' was not found", e.TaskId)

}
