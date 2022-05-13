package model

import (
	"time"
)

type Task struct {
	ID          int64     `db:"id" json:"id,omitempty"`
	Title       string    `db:"title" json:"title,omitempty"`
	Status      string    `db:"status" json:"status,omitempty"`
	Created     time.Time `db:"created" json:"created"`
	Description string    `db:"description" json:"description,omitempty"`
	Assignee    int       `db:"assignee" json:"assignee,omitempty"`
}
