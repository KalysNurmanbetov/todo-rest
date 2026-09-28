package repositories

import (
	"context"
	"errors"
	"fmt"
	"todo-rest/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewTaskRepositoryPostgres(p *pgxpool.Pool) *TaskRepositoryPostgres {
	return &TaskRepositoryPostgres{pool: p}
}

func (r *TaskRepositoryPostgres) Insert(t *domain.Task) (*domain.Task, error) {
	op := "infra.postgres.repositories.TaskRepositoryPostgres.Insert"
	ctx := context.Background()

	m := fromDomainToModel(t)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO tasks (id, title, completed, created_at, completed_at)
		VALUES($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET title = EXCLUDED.title,
			completed = EXCLUDED.completed,
			completed_at = EXCLUDED.completed_at
	`, m.Id, m.Title, m.Completed, m.CreatedAt, m.CompletedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return t, nil
}

func (r *TaskRepositoryPostgres) RemoveById(id domain.TaskId) error {
	op := "infra.postgres.repositories.TaskRepositoryPostgres.RemoveById"
	ctx := context.Background()

	tag, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id=$1`, id.Value())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}
	return nil
}

func (r *TaskRepositoryPostgres) FindById(id domain.TaskId) (*domain.Task, error) {
	op := "infra.postgres.repositories.TaskRepositoryPostgres.FindById"
	ctx := context.Background()

	var m TaskModel
	err := r.pool.QueryRow(ctx, `
		SELECT id, title, completed, created_at, completed_at
		FROM tasks
		WHERE id=$1
	`, id.Value()).Scan(&m.Id, &m.Title, &m.Completed, &m.CreatedAt, &m.CompletedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrTaskNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return m.toDomain()
}

func (r *TaskRepositoryPostgres) FindAll() ([]*domain.Task, error) {
	op := "infra.postgres.repositories.TaskRepositoryPostgres.FindAll"
	ctx := context.Background()

	rows, err := r.pool.Query(ctx, `
		SELECT id, title, completed, created_at, completed_at FROM tasks
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return prepareTaskList(rows, op)
}

func (r *TaskRepositoryPostgres) FindAllCompleted() ([]*domain.Task, error) {
	op := "infra.postgres.repositories.TaskRepositoryPostgres.FindAllCompleted"
	ctx := context.Background()

	rows, err := r.pool.Query(ctx, `
		SELECT id, title, completed, created_at, completed_at FROM tasks WHERE completed=$1
	`, true)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return prepareTaskList(rows, op)
}

func (r *TaskRepositoryPostgres) FindAllUncompleted() ([]*domain.Task, error) {
	op := "infra.postgres.repositories.TaskRepositoryPostgres.FindAllUncompleted"
	ctx := context.Background()

	rows, err := r.pool.Query(ctx, `
		SELECT id, title, completed, created_at, completed_at FROM tasks WHERE completed=$1
	`, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return prepareTaskList(rows, op)
}

// Helpers
func prepareTaskList(rows pgx.Rows, op string) ([]*domain.Task, error) {
	result := []*domain.Task{}
	for rows.Next() {
		var m TaskModel
		if err := rows.Scan(&m.Id, &m.Title, &m.Completed, &m.CreatedAt, &m.CompletedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		task, err := m.toDomain()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		result = append(result, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func fromDomainToModel(d *domain.Task) TaskModel {
	copy := d
	return TaskModel{
		Id:          copy.Id().Value(),
		Title:       copy.Title().Value(),
		Completed:   copy.Completed(),
		CreatedAt:   copy.CreatedAt(),
		CompletedAt: copy.CompletedAt(),
	}
}

func (m TaskModel) toDomain() (*domain.Task, error) {
	id, err := domain.TaskIdFrom(m.Id)
	if err != nil {
		return nil, fmt.Errorf("corrupted task id %q in storage: %w", m.Id, err)
	}
	title, err := domain.NewTaskTitle(m.Title)
	if err != nil {
		return nil, fmt.Errorf("corrupted task title for id %q in storage: %w", m.Id, err)
	}
	return domain.RehydrateTask(id, title, m.Completed, m.CreatedAt, m.CompletedAt), nil

}
