package infra

import (
	"fmt"
	"time"
	"todo-rest/domain"
)

type taskModel struct {
	Id          string
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

type TaskRepostoryMap struct {
	tasks map[string]taskModel
}

func NewTaskRepository() *TaskRepostoryMap {
	return &TaskRepostoryMap{tasks: map[string]taskModel{}}
}

func (r *TaskRepostoryMap) Insert(t *domain.Task) (*domain.Task, error) {
	r.tasks[t.Id().Value()] = fromDomainToModel(t)
	return t, nil
}

func (r *TaskRepostoryMap) RemoveById(id domain.TaskId) error {
	if _, ok := r.tasks[id.Value()]; ok == true {
		delete(r.tasks, id.Value())
		return nil
	} else {
		return domain.ErrTaskNotFound
	}
}

func (r *TaskRepostoryMap) FindById(id domain.TaskId) (*domain.Task, error) {
	if t, ok := r.tasks[id.Value()]; ok == true {
		return t.toDomain()
	} else {
		return nil, domain.ErrTaskNotFound
	}
}

func (r *TaskRepostoryMap) FindAll() ([]*domain.Task, error) {
	return r.getAllTasks(func(t taskModel) bool {
		return true
	})
}

func (r *TaskRepostoryMap) FindAllCompleted() ([]*domain.Task, error) {
	return r.getAllTasks(func(t taskModel) bool {
		return t.Completed
	})

}

func (r *TaskRepostoryMap) FindAllUncompleted() ([]*domain.Task, error) {
	return r.getAllTasks(func(t taskModel) bool {
		return t.Completed == false
	})
}

// Helpers
func fromDomainToModel(d *domain.Task) taskModel {
	copy := d
	return taskModel{
		Id:          copy.Id().Value(),
		Title:       copy.Title().Value(),
		Completed:   copy.Completed(),
		CreatedAt:   copy.CreatedAt(),
		CompletedAt: copy.CompletedAt(),
	}
}

func (m taskModel) toDomain() (*domain.Task, error) {
	id, err := domain.TaskIdFromString(m.Id)
	if err != nil {
		return nil, fmt.Errorf("corrupted task id %q in storage: %w", m.Id, err)
	}
	title, err := domain.NewTaskTitle(m.Title)
	if err != nil {
		return nil, fmt.Errorf("corrupted task title for id %q in storage: %w", m.Id, err)
	}
	return domain.RehydrateTask(id, title, m.Completed, m.CreatedAt, m.CompletedAt), nil

}

func (r *TaskRepostoryMap) getAllTasks(predicate func(t taskModel) bool) ([]*domain.Task, error) {
	result := []*domain.Task{}
	for _, t := range r.tasks {
		if predicate(t) {
			if task, err := t.toDomain(); err != nil {
				return nil, err
			} else {
				result = append(result, task)
			}
		}
	}
	return result, nil
}
