package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const os = "infra.config.postgres.helper"

func errSource(funcName string) string {
	return fmt.Sprintf("%s.%s", os, funcName)
}

func OpenPool(ctx context.Context, dns string) (*pgxpool.Pool, error) {
	funcName := "OpenPool"
	pool, err := pgxpool.New(ctx, dns)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", errSource(funcName), err)
	}
	return pool, nil
}

func BootstrapSchema(ctx context.Context, p *pgxpool.Pool) error {
	funcName := "BootstrapSchema"

	if _, err := p.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS tasks(
		id text PRIMARY KEY,
		title text NOT NULL,
		completed boolean NOT NULL DEFAULT false,
		created_at timestamptz NOT NULL,
		completed_at timestamptz
	)
	`); err != nil {
		return fmt.Errorf("%s:%w", errSource(funcName), err)
	}
	return nil
}
