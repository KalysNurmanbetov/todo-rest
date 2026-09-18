package infra

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"todo-rest/application"
	"todo-rest/domain"
)

type TaskHandler struct {
	service *application.TaskService
}

func NewTaskHandler(service *application.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tasks/{id}", h.GetTask)
	mux.HandleFunc("POST /tasks", h.CreateTask)
	mux.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)
	mux.HandleFunc("PATCH /tasks/{id}", h.ChangeTaskTitle)
	mux.HandleFunc("POST /tasks/{id}/complete", h.CompleteTask)
	mux.HandleFunc("POST /tasks/{id}/uncomplete", h.UncompleteTask)
	mux.HandleFunc("GET /tasks", h.GetAllTasks)
}

//Routes

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if task, err := h.service.GetTask(id); err != nil {
		err = fmt.Errorf("task %q: %w", id, err)
		responseWithError(w, err, http.StatusNotFound)
	} else {
		responseWithJsonBody(w, ToTaskDtoFrom(task), http.StatusOK)
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var dto TaskTitleDto

	if err := decodeAndValidateDto(r.Body, &dto); err != nil {
		responseWithError(w, err, http.StatusBadRequest)
		return
	}

	if task, err := h.service.InsertNewTask(dto.Title); err != nil {
		err = fmt.Errorf("title %q: %w", dto.Title, err)
		responseWithError(w, err, http.StatusConflict)
	} else {
		responseWithJsonBody(w, CreatedTaskDto{Id: task.Id().Value(), Title: task.Title().Value()}, http.StatusCreated)
	}
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.DeleteTask(id); err != nil {
		responseWithError(w, err, http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *TaskHandler) ChangeTaskTitle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var dto TaskTitleDto
	if err := decodeAndValidateDto(r.Body, &dto); err != nil {
		responseWithError(w, err, http.StatusBadRequest)
		return
	}

	if task, err := h.service.ChangeTaskTitle(id, dto.Title); err != nil {
		err = fmt.Errorf("task %q: %w", id, err)
		var code int
		if errors.Is(err, domain.ErrTaskNotFound) {
			code = http.StatusNotFound
		} else {
			code = http.StatusConflict
		}

		responseWithError(w, err, code)
	} else {
		responseWithJsonBody(w, ToTaskDtoFrom(task), http.StatusOK)
	}
}

func (h *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if task, err := h.service.CompleteTask(id); err != nil {
		err = fmt.Errorf("task %q: %w", id, err)
		var code int
		if errors.Is(err, domain.ErrTaskNotFound) {
			code = http.StatusNotFound
		} else {
			code = http.StatusConflict
		}

		responseWithError(w, err, code)
	} else {
		responseWithJsonBody(w, ComletedTaskDto{Id: task.Id().Value(), Completed: task.Completed(), CompletedAt: *task.CompletedAt()}, http.StatusOK)
	}
}

func (h *TaskHandler) UncompleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if task, err := h.service.UncompleteTask(id); err != nil {
		err = fmt.Errorf("task %q: %w", id, err)
		var code int
		if errors.Is(err, domain.ErrTaskNotFound) {
			code = http.StatusNotFound
		} else {
			code = http.StatusConflict
		}
		responseWithError(w, err, code)
	} else {
		responseWithJsonBody(w, UncomletedTaskDto{Id: task.Id().Value(), Completed: task.Completed()}, http.StatusOK)
	}
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []*domain.Task
	completed := r.URL.Query().Get("completed")

	var f application.CompletionFilter
	switch completed {
	case "true":
		f = application.AllCompleted
	case "false":
		f = application.AllUncompleted
	default:
		f = application.All
	}

	tasks = h.service.GetAllTasks(f)

	taskDtos := []TaskDto{}
	for _, t := range tasks {
		taskDtos = append(taskDtos, ToTaskDtoFrom(t))
	}

	responseWithJsonBody(w, taskDtos, http.StatusOK)
}

//Helpers

func decodeAndValidateDto(r io.Reader, dto Validatable) error {
	if err := json.NewDecoder(r).Decode(dto); err != nil {
		return err
	}
	if err := dto.Validate(); err != nil {
		return err
	}
	return nil
}

func responseWithError(w http.ResponseWriter, err error, code int) {
	responseWithJsonBody(w, NewErrorDto(err), code)
}

func responseWithJsonBody(w http.ResponseWriter, jsonBody any, code int) {
	b, err := json.MarshalIndent(jsonBody, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}
