package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          int64     `db:"id" json:"id,omitempty"`
	PublicID    uuid.UUID `db:"public_id" json:"public_id"`
	Key         int       `db:"key" json:"key"`
	Title       string    `db:"title" json:"title"`
	Status      string    `db:"status" json:"status"`
	Created     time.Time `db:"created" json:"created"`
	Description string    `db:"description" json:"description"`
	Assignee    int       `db:"assignee" json:"assignee"`
}

type NewTaskEvent struct {
	ID          int64     `json:"id,omitempty"`
	PublicID    uuid.UUID `json:"public_id"`
	Key         string    `json:"key"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	Created     time.Time `json:"created"`
	Description string    `json:"description"`
	Assignee    int       `json:"assignee"`
}

type TaskCompletedEvent struct {
	ID       uuid.UUID `json:"id"`
	Assignee int       `json:"assignee"`
}
