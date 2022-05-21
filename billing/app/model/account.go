package model

import (
	"github.com/guregu/null"
	"github.com/shopspring/decimal"
)

type AccountSmall struct {
	ID        int             `db:"id"`
	Balance   decimal.Decimal `db:"balance"`
	UserEmail string          `db:"user_email"`
}

type Account struct {
	Balance decimal.Decimal `db:"balance" json:"balance"`
	History []HistoryItem   `db:"history" json:"history"`
}

type HistoryItem struct {
	Type            string              `db:"type" json:"type"`
	BalanceChange   decimal.Decimal     `db:"balance_change" json:"balance_change"`
	TaskTitle       null.String         `db:"task_title" json:"task_title"`
	TaskDescription null.String         `db:"task_description" json:"task_description"`
	TaskReward      decimal.NullDecimal `db:"task_reward" json:"task_reward"`
}
