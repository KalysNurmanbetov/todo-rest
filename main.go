package main

import (
	"fmt"
	"net/http"
	"todo-rest/application"
	"todo-rest/infra"
	"todo-rest/infra/config"
)

func main() {
	cfg := config.MustLoad()

	mux := http.NewServeMux()

	taskRepository := infra.NewTaskRepository()
	taskService := application.NewTaskService(taskRepository)
	taskHandler := infra.NewTaskHandler(taskService)

	taskHandler.RegisterRoutes(mux)

	if err := http.ListenAndServe(cfg.Address, mux); err != nil {
		fmt.Println("Server start failed: ", err)
	}
}
