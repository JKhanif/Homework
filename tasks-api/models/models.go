package models

import (
	"time"
)

type CreateTask struct {
	Name        string `json:"Название"`
	Description string `json:"Описание,omitempty"`
}

type TaskResponse struct {
	Id          int       `json:"N"`
	Name        string    `json:"Название"`
	Description string    `json:"Описание,omitempty"`
	CreatedAt   time.Time `json:"Когда"`
}
