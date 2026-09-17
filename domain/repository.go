package domain

import (
	"sync"
)

type TaskRepostory struct {
	tasks map[string]*Task
	mtx   sync.RWMutex
}

func NewTaskRepository() *TaskRepostory {
	return &TaskRepostory{tasks: map[string]*Task{}}
}

func (r *TaskRepostory) Insert(t *Task) *Task {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.tasks[t.Id] = t
	return t
}

func (r *TaskRepostory) Update(id string, fn func(*Task) error) (*Task, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	t, ok := r.tasks[id]
	if !ok {
		return nil, ErrTaskNotFound
	}
	copy := *t

	if err := fn(&copy); err != nil {
		return nil, err
	}
	r.tasks[id] = &copy
	return &copy, nil
}

func (r *TaskRepostory) RemoveById(id string) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.tasks[id]; ok == true {
		delete(r.tasks, id)
		return true
	} else {
		return false
	}
}

func (r *TaskRepostory) FindById(id string) (*Task, bool) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	t, ok := r.tasks[id]
	return t, ok
}

func (r *TaskRepostory) FindAll() []*Task {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	result := []*Task{}
	for _, t := range r.tasks {
		result = append(result, t)
	}
	return result
}

func (r *TaskRepostory) FindAllCompleted() []*Task {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	result := []*Task{}
	for _, t := range r.tasks {
		if t.Completed {
			result = append(result, t)
		}
	}
	return result
}

func (r *TaskRepostory) FindAllUncompleted() []*Task {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	result := []*Task{}
	for _, t := range r.tasks {
		if t.Completed == false {
			result = append(result, t)
		}
	}
	return result
}
