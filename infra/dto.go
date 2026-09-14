package infra

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Validatable interface {
	Validate() error
}

type CreateTaskDto struct {
	Title string `json:"title"`
}

func (d CreateTaskDto) Validate() error {
	if strings.Trim(d.Title, " ") == "" {
		return errors.New("title is empty")
	}
	return nil
}

type TaskCreteadDto struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type TaskDto struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type ErrorDto struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func NewErrorDto(msg string) *ErrorDto {
	return &ErrorDto{Message: msg, Time: time.Now()}
}

func (d ErrorDto) ToString() string {
	if b, err := json.MarshalIndent(d, "", "    "); err != nil {
		panic(err)
	} else {
		return string(b)
	}
}
