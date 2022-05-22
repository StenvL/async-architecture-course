package repository

import (
	"time"

	"github.com/shopspring/decimal"
)

func (r Repository) AddBalanceChange(accountID int, balanceChanging decimal.Decimal, timestamp time.Time) error {
	query := "insert into balance_changes (account_id, balance_changing, timestamp) values ($1, $2, $3)"
	_, err := r.db.Exec(query, accountID, balanceChanging, timestamp)
	return err
}

func (r Repository) GetDailyIncome(day time.Time) (decimal.Decimal, error) {
	var income decimal.Decimal
	query := `
		select -sum(balance_changing) as income
		from balance_changes	
		where "timestamp" between $1 and $2`

	err := r.db.Get(&income, query, day, day.Add(1*time.Hour*24))
	return income, err
}

func (r Repository) GetDailyAccountsWithNegativeBalance(day time.Time) (int, error) {
	var count int
	query := `
		select count(*) from (
			select 1
			from balance_changes
			where "timestamp" >= $1 and "timestamp" < $2
			group by account_id
			having sum(balance_changing) < 0
		) t1`

	err := r.db.Get(&count, query, day, day.Add(1*time.Hour*24))
	return count, err
}

func (r Repository) GetMostExpensiveTask(from, to time.Time) (decimal.Decimal, error) {
	var reward decimal.Decimal
	query := `select max(reward) as reward from tasks where "timestamp" between $1 and $2`
	err := r.db.Get(&reward, query, from, to)

	return reward, err
}
