package main

import (
	"fmt"
	"net/http"
	"todo-rest/domain"
	"todo-rest/infra"
)

func main() {
	taskRepository := domain.NewTaskRepository()
	taskController := infra.NewTaskController(taskRepository, &infra.IdGenerator{})
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", taskController.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", taskController.GetTask)
	mux.HandleFunc("GET /tasks", taskController.GetAllTasks)
	mux.HandleFunc("POST /tasks/{id}/complete", taskController.CompleteTask) //could also do PATCH with request body {complete: bool}
	mux.HandleFunc("POST /tasks/{id}/uncomplete", taskController.UncompleteTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskController.DeleteTask)

	if err := http.ListenAndServe(":9091", mux); err != nil {
		fmt.Println("Server start failed: ", err)
	}
}
