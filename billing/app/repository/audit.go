package repository

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"

	"github.com/StenvL/async-architecture-course/billing/app/model"
)

func (r Repository) GetBalanceChanges(date time.Time) ([]decimal.Decimal, error) {
	var res []decimal.Decimal

	query := `
		select balance_change
		from accounts_audit_log
		where timestamp >= CURRENT_DATE and timestamp <= CURRENT_DATE + interval '1 day'`
	err := r.db.Select(&res, query)

	return res, err
}

func (r Repository) AddToAuditLog(record model.AuditLogRecord) error {
	query := `
		insert into accounts_audit_log (account_id, task_id, type, balance_change) 
		values (:account_id, :task_id, :type, :balance_change)`
	_, err := r.db.NamedExec(query, record)
	return err
}

func (r Repository) AddToAuditLogTx(tx *sql.Tx, record model.AuditLogRecord) error {
	query := `
		insert into accounts_audit_log (account_id, task_id, type, balance_change) 
		values ($1, $2, $3, $4)`
	_, err := tx.Exec(query, record.AccountID, record.TaskID, record.Type, record.BalanceChange)
	return err
}
