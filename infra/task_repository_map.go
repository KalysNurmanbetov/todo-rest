package infra

import (
	"todo-rest/domain/task"
)

type TaskRepostoryMap struct {
	tasks map[string]*task.Task
}

func NewTaskRepository() *TaskRepostoryMap {
	return &TaskRepostoryMap{tasks: map[string]*task.Task{}}
}

func (r *TaskRepostoryMap) Insert(t *task.Task) (*task.Task, error) {
	r.tasks[t.Id().Value()] = t
	return t, nil
}

func (r *TaskRepostoryMap) RemoveById(id string) error {
	if _, ok := r.tasks[id]; ok == true {
		delete(r.tasks, id)
		return nil
	} else {
		return task.ErrTaskNotFound
	}
}

func (r *TaskRepostoryMap) FindById(id string) (*task.Task, bool) {
	t, ok := r.tasks[id]
	return t, ok
}

func (r *TaskRepostoryMap) FindAll() []*task.Task {
	result := []*task.Task{}
	for _, t := range r.tasks {
		result = append(result, t)
	}
	return result
}

func (r *TaskRepostoryMap) FindAllCompleted() []*task.Task {
	result := []*task.Task{}
	for _, t := range r.tasks {
		if t.Completed() {
			result = append(result, t)
		}
	}
	return result
}

func (r *TaskRepostoryMap) FindAllUncompleted() []*task.Task {
	result := []*task.Task{}
	for _, t := range r.tasks {
		if t.Completed() == false {
			result = append(result, t)
		}
	}
	return result
}
