package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	TaskAssignedRecordType  = "task_assigned"
	TaskCompletedRecordType = "task_completed"
	PaymentMadeRecordType   = "payment_made"
)

type AuditLogRecord struct {
	AccountID     int             `db:"account_id"`
	TaskID        uuid.NullUUID   `db:"task_id"`
	Type          string          `db:"type"`
	BalanceChange decimal.Decimal `db:"balance_change"`
}
