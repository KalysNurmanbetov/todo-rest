package infra

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
	"todo-rest/domain"
)

var ErrTitleIsNotProvided = errors.New("title must be defined")

type Validatable interface {
	Validate() error
}

type TaskTitleDto struct {
	Title string `json:"title"`
}

func (d TaskTitleDto) Validate() error {
	if strings.Trim(d.Title, " ") == "" {
		return ErrTitleIsNotProvided
	}
	return nil
}

type CreatedTaskDto struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type TaskDto struct {
	Id          string     `json:"id"`
	Title       string     `json:"title"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

func ToTaskDtoFrom(t *domain.Task) TaskDto {
	return TaskDto{
		Id:          t.Id().Value(),
		Title:       t.Title().Value(),
		Completed:   t.Completed(),
		CreatedAt:   t.CreatedAt(),
		CompletedAt: t.CompletedAt(),
	}
}

type ErrorDto struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func NewErrorDto(err error) *ErrorDto {
	return &ErrorDto{Message: err.Error(), Time: time.Now()}
}

func (d ErrorDto) ToString() string {
	if b, err := json.MarshalIndent(d, "", "    "); err != nil {
		panic(err)
	} else {
		return string(b)
	}
}
