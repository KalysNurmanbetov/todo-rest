package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"todo-rest/application"
	"todo-rest/infra"
	"todo-rest/infra/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()

	mux := http.NewServeMux()

	pool, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("connection pool is not created: %s", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping is not succesed, %s", err)
	}

	_, err = pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS tasks(
		id text PRIMARY KEY,
		title text NOT NULL,
		completed boolean NOT NULL DEFAULT false,
		created_at timestamptz NOT NULL,
		completed_at timestamptz
	)
	`)
	if err != nil {
		log.Fatalf("not created, %s", err)
	}

	taskRepository := infra.NewTaskRepository()
	taskService := application.NewTaskService(taskRepository)
	taskHandler := infra.NewTaskHandler(taskService)

	taskHandler.RegisterRoutes(mux)

	if err := http.ListenAndServe(cfg.Address, mux); err != nil {
		fmt.Println("Server start failed: ", err)
	}
}
