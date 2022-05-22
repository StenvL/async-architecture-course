package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (r Repository) TaskEstimated(taskID uuid.UUID, reward decimal.Decimal, timestamp time.Time) error {
	_, err := r.db.Exec("insert into tasks (id, reward, timestamp) values ($1, $2, $3)", taskID, reward, timestamp)
	return err
}
