package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/shopspring/decimal"
)

type Task struct {
	ID          uuid.UUID       `db:"id" json:"id,omitempty"`
	Title       string          `db:"title" json:"title,omitempty"`
	Key         string          `db:"key" json:"key,omitempty"`
	Status      string          `db:"status" json:"status,omitempty"`
	Created     time.Time       `db:"created" json:"created"`
	Description string          `db:"description" json:"description,omitempty"`
	Assignee    int             `db:"assignee" json:"assignee,omitempty"`
	Cost        decimal.Decimal `db:"cost" json:"cost,omitempty"`
	Reward      decimal.Decimal `db:"reward" json:"reward,omitempty"`
}
