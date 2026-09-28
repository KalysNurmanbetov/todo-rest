package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"todo-rest/application"
	"todo-rest/infra"
	"todo-rest/infra/config"
	"todo-rest/infra/postgres"
	"todo-rest/infra/postgres/repositories"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()

	mux := http.NewServeMux()

	pool, err := postgres.OpenPool(ctx, cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("pool is not created: %s", err)
	}
	defer pool.Close()

	err = postgres.BootstrapSchema(ctx, pool)
	if err != nil {
		log.Fatalf("bootstrap schema failed: %s", err)
	}

	taskRepository := repositories.NewTaskRepositoryPostgres(pool)
	taskService := application.NewTaskService(taskRepository)
	taskHandler := infra.NewTaskHandler(taskService)

	taskHandler.RegisterRoutes(mux)

	if err := http.ListenAndServe(cfg.Address, mux); err != nil {
		fmt.Println("Server start failed: ", err)
	}
}
