package repository

import (
	"database/sql"

	"github.com/StenvL/async-architecture-course/billing/app/model"

	"github.com/shopspring/decimal"
)

func (r Repository) GetUserAccount(userID int) (model.Account, error) {
	var res model.Account

	if err := r.db.Get(&res, "select balance from accounts where user_id = $1", userID); err != nil {
		return model.Account{}, err
	}

	history, err := r.getAccountHistory(userID)
	if err != nil {
		return model.Account{}, err
	}
	res.History = history

	return res, nil
}

func (r Repository) CreateAccount(userID int) error {
	_, err := r.db.Exec("insert into accounts (user_id, balance) values ($1, 0)", userID)
	return err
}

func (r Repository) ChangeAccountBalanceTx(tx *sql.Tx, accountID int, balanceChange decimal.Decimal) error {
	_, err := tx.Exec("update accounts set balance = balance + $1 where id = $2", balanceChange, accountID)
	return err
}

func (r Repository) getPositiveBalancedAccounts() ([]model.AccountSmall, error) {
	var res []model.AccountSmall

	query := `
		select a.id, a.balance, u.email as user_email
		from accounts a join users u on a.user_id = u.id 
		where balance > 0`
	if err := r.db.Select(&res, query); err != nil {
		return nil, err
	}

	return res, nil
}

func (r Repository) getAccountHistory(userID int) ([]model.HistoryItem, error) {
	var res []model.HistoryItem

	query := `
		select aul.type, aul.balance_change, coalesce('[' || t.key || '] - ' || t.title, t.title) as task_title, t.description as task_description, t.reward as task_reward
		from accounts_audit_log aul
			join accounts a on aul.account_id = a.id
			left join tasks t on aul.task_id = t.id
		where a.user_id = $1
		order by aul.id desc`
	if err := r.db.Select(&res, query, userID); err != nil {
		return nil, err
	}

	return res, nil
}
