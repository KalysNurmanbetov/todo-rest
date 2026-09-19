package application

import (
	"sync"
	"todo-rest/domain"
)

type CompletionFilter string

const (
	All            CompletionFilter = "all"
	AllCompleted   CompletionFilter = "completed"
	AllUncompleted CompletionFilter = "uncompleted"
)

type TaskService struct {
	taskRepository domain.TaskRepository
	mtx            sync.RWMutex
}

func NewTaskService(repository domain.TaskRepository) *TaskService {
	return &TaskService{taskRepository: repository}
}

func (ts *TaskService) InsertNewTask(title string) (*domain.Task, error) {
	taskTitle, err := domain.NewTaskTitle(title)
	if err != nil {
		return nil, err
	}

	newTask := domain.NewTask(domain.NewTaskId(), taskTitle)

	ts.mtx.Lock()
	defer ts.mtx.Unlock()
	return ts.taskRepository.Insert(newTask)
}

func (ts *TaskService) GetTask(id string) (*domain.Task, error) {

	taskId, err := domain.TaskIdFrom(id)

	if err != nil {
		return nil, err
	}

	ts.mtx.RLock()
	defer ts.mtx.RUnlock()
	return ts.taskRepository.FindById(taskId)
}

func (ts *TaskService) DeleteTask(id string) error {
	taskId, err := domain.TaskIdFrom(id)

	if err != nil {
		return err
	}

	ts.mtx.Lock()
	defer ts.mtx.Unlock()
	return ts.taskRepository.RemoveById(taskId)
}

func (ts *TaskService) GetAllTasks(f CompletionFilter) ([]*domain.Task, error) {
	var tasks []*domain.Task
	var err error

	ts.mtx.RLock()
	defer ts.mtx.RUnlock()

	switch f {
	case AllCompleted:
		tasks, err = ts.taskRepository.FindAllCompleted()
	case AllUncompleted:
		tasks, err = ts.taskRepository.FindAllUncompleted()
	case All:
	default:
		tasks, err = ts.taskRepository.FindAll()
	}

	return tasks, err
}

func (ts *TaskService) ChangeTaskTitle(id, title string) (*domain.Task, error) {
	taskId, err := domain.TaskIdFrom(id)
	if err != nil {
		return nil, err
	}

	newTitle, err := domain.NewTaskTitle(title)
	if err != nil {
		return nil, err
	}

	ts.mtx.Lock()
	defer ts.mtx.Unlock()

	task, err := ts.taskRepository.FindById(taskId)
	if err != nil {
		return nil, err
	}

	if err := task.ChangeTitle(newTitle); err != nil {
		return nil, err
	}

	return ts.taskRepository.Insert(task)
}

func (ts *TaskService) CompleteTask(id string) (*domain.Task, error) {
	taskId, err := domain.TaskIdFrom(id)
	if err != nil {
		return nil, err
	}

	ts.mtx.Lock()
	defer ts.mtx.Unlock()

	task, err := ts.taskRepository.FindById(taskId)
	if err != nil {
		return nil, err
	}

	if err := task.Complete(); err != nil {
		return nil, err
	}

	return ts.taskRepository.Insert(task)
}

func (ts *TaskService) UncompleteTask(id string) (*domain.Task, error) {
	taskId, err := domain.TaskIdFrom(id)
	if err != nil {
		return nil, err
	}

	ts.mtx.Lock()
	defer ts.mtx.Unlock()

	task, err := ts.taskRepository.FindById(taskId)
	if err != nil {
		return nil, err
	}

	if err := task.Uncomplete(); err != nil {
		return nil, err
	}

	return ts.taskRepository.Insert(task)
}
