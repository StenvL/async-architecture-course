package api

import "github.com/shopspring/decimal"

type dailyIncomeResponse struct {
	Sum   decimal.Decimal `json:"sum"`
	Count int             `json:"count"`
}

type mostExpensiveTaskResponse struct {
	Sum decimal.Decimal `json:"sum"`
}
