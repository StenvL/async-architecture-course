package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"

	"github.com/StenvL/async-architecture-course/billing/app/model"
	"github.com/shopspring/decimal"
)

func (r Repository) GetTaskCost(taskID uuid.UUID) (decimal.Decimal, error) {
	var taskCost decimal.Decimal
	if err := r.db.Get(&taskCost, "select cost from tasks where id = $1", taskID); err != nil {
		return decimal.Decimal{}, err
	}

	return taskCost, nil
}

func (r Repository) CreateTask(task model.Task) (int, error) {
	accountID, err := r.GetAccountID(task.Assignee)
	if err != nil {
		return 0, err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	query := `
		insert into tasks (id, title, key, status, created, description, assignee, cost, reward)
		values (:id, :title, :key, :status, :created, :description, :assignee, :cost, :reward)`
	query, args, err := sqlx.Named(query, task)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	query = sqlx.Rebind(sqlx.BindType("postgres"), query)
	_, err = tx.Exec(query, args...)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	if err = r.assignTaskTx(tx, accountID, task.ID, task.Cost.Neg()); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return accountID, tx.Commit()
}

func (r Repository) CompleteTask(taskID uuid.UUID, userID int) (int, decimal.Decimal, error) {
	var taskReward decimal.Decimal
	if err := r.db.Get(&taskReward, "select reward from tasks where id = $1", taskID); err != nil {
		return 0, decimal.Decimal{}, err
	}

	var accountID int
	if err := r.db.Get(&accountID, "select id from accounts where user_id = $1", userID); err != nil {
		return 0, decimal.Decimal{}, err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, decimal.Decimal{}, err
	}

	_, err = tx.Exec("update tasks set status = 'resolved' where id = $1", taskID)
	if err != nil {
		_ = tx.Rollback()
		return 0, decimal.Decimal{}, err
	}

	if err = r.AddToAuditLogTx(tx, model.AuditLogRecord{
		AccountID:     accountID,
		TaskID:        uuid.NullUUID{UUID: taskID, Valid: true},
		Type:          model.TaskCompletedRecordType,
		BalanceChange: taskReward,
	}); err != nil {
		_ = tx.Rollback()
		return 0, decimal.Decimal{}, err
	}

	if err = r.ChangeAccountBalanceTx(tx, accountID, taskReward); err != nil {
		_ = tx.Rollback()
		return 0, decimal.Decimal{}, err
	}

	return accountID, taskReward, tx.Commit()
}

func (r Repository) UpdateShuffledTasks(shuffledTasks map[uuid.UUID]int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for taskID, userID := range shuffledTasks {
		taskCost, err := r.GetTaskCost(taskID)
		if err != nil {
			return err
		}

		accountID, err := r.GetAccountID(userID)
		if err != nil {
			return err
		}

		if _, err = tx.Exec("update tasks set assignee = $1 where id = $2", userID, taskID); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err = r.assignTaskTx(tx, accountID, taskID, taskCost.Neg()); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r Repository) assignTaskTx(tx *sql.Tx, accountID int, taskID uuid.UUID, balanceChange decimal.Decimal) error {
	if err := r.AddToAuditLogTx(tx, model.AuditLogRecord{
		AccountID:     accountID,
		TaskID:        uuid.NullUUID{UUID: taskID, Valid: true},
		Type:          model.TaskAssignedRecordType,
		BalanceChange: balanceChange,
	}); err != nil {
		return err
	}

	if err := r.ChangeAccountBalanceTx(tx, accountID, balanceChange); err != nil {
		return err
	}

	return nil
}
