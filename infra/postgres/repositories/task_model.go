package repositories

import "time"

type TaskModel struct {
	Id          string
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}
