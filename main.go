package main

import (
	"fmt"
	"net/http"
	"todo-rest/application"
	"todo-rest/infra"
)

func main() {
	mux := http.NewServeMux()

	taskRepository := infra.NewTaskRepository()
	taskService := application.NewTaskService(taskRepository)
	taskHandler := infra.NewTaskHandler(taskService)

	taskHandler.RegisterRoutes(mux)

	if err := http.ListenAndServe(":9091", mux); err != nil {
		fmt.Println("Server start failed: ", err)
	}
}
