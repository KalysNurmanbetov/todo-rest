package infra

import (
	"todo-rest/domain"
)

type TaskRepostoryMap struct {
	tasks map[string]*domain.Task
}

func NewTaskRepository() *TaskRepostoryMap {
	return &TaskRepostoryMap{tasks: map[string]*domain.Task{}}
}

func (r *TaskRepostoryMap) Insert(t *domain.Task) (*domain.Task, error) {
	r.tasks[t.Id().Value()] = t
	return t, nil
}

func (r *TaskRepostoryMap) RemoveById(id string) error {
	if _, ok := r.tasks[id]; ok == true {
		delete(r.tasks, id)
		return nil
	} else {
		return domain.ErrTaskNotFound
	}
}

func (r *TaskRepostoryMap) FindById(id string) (*domain.Task, bool) {
	t, ok := r.tasks[id]
	return t, ok
}

func (r *TaskRepostoryMap) FindAll() []*domain.Task {
	result := []*domain.Task{}
	for _, t := range r.tasks {
		result = append(result, t)
	}
	return result
}

func (r *TaskRepostoryMap) FindAllCompleted() []*domain.Task {
	result := []*domain.Task{}
	for _, t := range r.tasks {
		if t.Completed() {
			result = append(result, t)
		}
	}
	return result
}

func (r *TaskRepostoryMap) FindAllUncompleted() []*domain.Task {
	result := []*domain.Task{}
	for _, t := range r.tasks {
		if t.Completed() == false {
			result = append(result, t)
		}
	}
	return result
}
