package models

import (
	"time"
)

type Todo struct {
	ID        int       `json:"id,omitempty"`
	TaskName  string    `json:"task_name" validate:"required,min=5"`
	Completed bool      `json:"completed"`
	DueDate   time.Time `json:"due_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
