package task

import "todo-rest/domain/task"

type TaskService struct {
	taskRepository task.TaskRepostory
}

func NewTaskService(repository task.TaskRepostory) *TaskService {
	return &TaskService{taskRepository: repository}
}

func (ts *TaskService) InsertNewTask(title string) (*task.Task, error) {
	taskTitle, err := task.NewTaskTitle(title)
	if err != nil {
		return nil, err
	}

	newTask, err := task.NewTask(task.NewTaskId(), taskTitle)
	if err != nil {
		return nil, err
	}

	return ts.taskRepository.Insert(newTask)
}

func (ts *TaskService) GetTask(id string) (*task.Task, error) {
	taskId, err := task.TaskIdFromString(id)

	if err != nil {
		return nil, err
	}

	return ts.taskRepository.FindById(taskId)
}

func (ts *TaskService) DeleteTask(id string) error {
	taskId, err := task.TaskIdFromString(id)

	if err != nil {
		return err
	}

	return ts.taskRepository.RemoveById(taskId)
}

// TODO: Add filtering
func (ts *TaskService) GetAllTasks() []*task.Task {
	return ts.taskRepository.FindAll()
}

func (ts *TaskService) ChangeTaskTitle(id string, title string) (*task.Task, error) {
	taskId, err := task.TaskIdFromString(id)
	if err != nil {
		return nil, err
	}

	newTitle, err := task.NewTaskTitle(title)
	if err != nil {
		return nil, err
	}

	task, err := ts.taskRepository.FindById(taskId)
	if err != nil {
		return nil, err
	}

	if err := task.ChangeTitle(newTitle); err != nil {
		return nil, err
	}

	return ts.taskRepository.Insert(task)
}

func (ts *TaskService) CompleteTask(id string) (*task.Task, error) {
	taskId, err := task.TaskIdFromString(id)
	if err != nil {
		return nil, err
	}

	task, err := ts.taskRepository.FindById(taskId)
	if err != nil {
		return nil, err
	}

	if err := task.Complete(); err != nil {
		return nil, err
	}

	return ts.taskRepository.Insert(task)
}

func (ts *TaskService) UncompleteTask(id string) (*task.Task, error) {
	taskId, err := task.TaskIdFromString(id)
	if err != nil {
		return nil, err
	}

	task, err := ts.taskRepository.FindById(taskId)
	if err != nil {
		return nil, err
	}

	if err := task.Uncomplete(); err != nil {
		return nil, err
	}

	return ts.taskRepository.Insert(task)
}
