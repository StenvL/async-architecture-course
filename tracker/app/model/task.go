package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          int64     `db:"id" json:"id,omitempty"`
	PublicID    uuid.UUID `db:"public_id" json:"public_id"`
	Title       string    `db:"title" json:"title"`
	Status      string    `db:"status" json:"status"`
	Created     time.Time `db:"created" json:"created"`
	Description string    `db:"description" json:"description"`
	Assignee    int       `db:"assignee" json:"assignee"`
}

type TaskCompletedEvent struct {
	ID       uuid.UUID `json:"id"`
	Assignee int       `json:"assignee"`
}
