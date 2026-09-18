package infra

import (
	"sync"
	"todo-rest/domain/task"
)

type TaskRepostoryMap struct {
	tasks map[string]*task.Task
	mtx   sync.RWMutex
}

func NewTaskRepository() *TaskRepostoryMap {
	return &TaskRepostoryMap{tasks: map[string]*task.Task{}}
}

func (r *TaskRepostoryMap) Insert(t *task.Task) (*task.Task, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.tasks[t.Id().Value()] = t
	return t, nil
}

func (r *TaskRepostoryMap) RemoveById(id string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.tasks[id]; ok == true {
		delete(r.tasks, id)
		return nil
	} else {
		return task.ErrTaskNotFound
	}
}

func (r *TaskRepostoryMap) FindById(id string) (*task.Task, bool) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	t, ok := r.tasks[id]
	return t, ok
}

func (r *TaskRepostoryMap) FindAll() []*task.Task {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	result := []*task.Task{}
	for _, t := range r.tasks {
		result = append(result, t)
	}
	return result
}

func (r *TaskRepostoryMap) FindAllCompleted() []*task.Task {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	result := []*task.Task{}
	for _, t := range r.tasks {
		if t.Completed() {
			result = append(result, t)
		}
	}
	return result
}

func (r *TaskRepostoryMap) FindAllUncompleted() []*task.Task {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	result := []*task.Task{}
	for _, t := range r.tasks {
		if t.Completed() == false {
			result = append(result, t)
		}
	}
	return result
}
