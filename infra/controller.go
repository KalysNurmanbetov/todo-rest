package infra

import (
	"errors"
	"fmt"
	"net/http"
	"todo-rest/domain"
)

type TaskController struct {
	taskRepository *domain.TaskRepostory
	idGenerator    *IdGenerator
}

func NewTaskController(taskRepository *domain.TaskRepostory, idGenerator *IdGenerator) *TaskController {
	return &TaskController{taskRepository: taskRepository, idGenerator: idGenerator}
}

func (c *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	var dto CreateTaskDto

	if err := DecodeAndValidateDto(r.Body, &dto); err != nil {
		CreateResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := domain.NewTask(c.idGenerator.Generate(), dto.Title)
	c.taskRepository.Insert(task)

	CreateResponseWithJsonBody(w, TaskCreteadDto{Id: task.Id, Title: task.Title}, http.StatusCreated)
}

func (c *TaskController) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		CreateResponseWithError(w, "id was not provided", http.StatusBadRequest)
		return
	}

	task, ok := c.taskRepository.FindById(id)
	if !ok {
		CreateResponseWithError(w, fmt.Sprintf("Task with id '%s' was not found", id), http.StatusNotFound)
		return
	}

	CreateResponseWithJsonBody(w, TaskDto{Id: task.Id, Title: task.Title, Completed: task.Completed}, http.StatusOK)
}

func (c *TaskController) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []domain.Task
	completed := r.URL.Query().Get("completed")

	switch completed {
	case "true":
		tasks = c.taskRepository.FindAllCompleted()
	case "false":
		tasks = c.taskRepository.FindAllUncompleted()
	default:
		tasks = c.taskRepository.FindAll()
	}

	taskDtos := []TaskDto{}

	for _, t := range tasks {
		taskDtos = append(taskDtos, TaskDto{Id: t.Id, Title: t.Title, Completed: t.Completed})
	}

	CreateResponseWithJsonBody(w, taskDtos, http.StatusOK)
}

func (c *TaskController) CompleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		CreateResponseWithError(w, "id was not provided", http.StatusBadRequest)
		return
	}

	task, err := c.taskRepository.Update(id, func(t *domain.Task) error {
		return t.Complete()
	})

	if err != nil {
		var notFoundErr domain.TaskNotFoundError
		var code int
		if errors.As(err, &notFoundErr) {
			code = http.StatusNotFound
		} else {
			code = http.StatusConflict
		}

		CreateResponseWithError(w, err.Error(), code)
		return
	}

	CreateResponseWithJsonBody(w, TaskDto{Id: task.Id, Title: task.Title, Completed: task.Completed}, http.StatusOK)

}

func (c *TaskController) UncompleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		CreateResponseWithError(w, "id was not provided", http.StatusBadRequest)
		return
	}

	task, err := c.taskRepository.Update(id, func(t *domain.Task) error {
		return t.Uncomplete()
	})

	if err != nil {
		var notFoundErr domain.TaskNotFoundError
		var code int
		if errors.As(err, &notFoundErr) {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}

		CreateResponseWithError(w, err.Error(), code)
		return
	}

	CreateResponseWithJsonBody(w, TaskDto{Id: task.Id, Title: task.Title, Completed: task.Completed}, http.StatusOK)

}

func (c *TaskController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		CreateResponseWithError(w, "id was not provided", http.StatusBadRequest)
		return
	}

	ok := c.taskRepository.RemoveById(id)

	if !ok {
		CreateResponseWithError(w, domain.TaskNotFoundError{TaskId: id}.Error(), http.StatusNotFound)
	}

	w.WriteHeader(http.StatusNoContent)
}
