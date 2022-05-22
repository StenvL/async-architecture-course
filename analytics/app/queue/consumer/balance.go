package consumer

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

func (c Client) ConsumeBalanceEvents() error {
	c.balanceChangedConsumer()

	return nil
}

func (c Client) balanceChangedConsumer() {
	handler := func(msg []byte) error {
		type eventPayload struct {
			AccountID       int             `json:"account_id"`
			BalanceChanging decimal.Decimal `json:"balance_changing"`
			Timestamp       time.Time       `json:"timestamp"`
		}
		msgStruct := struct {
			Data eventPayload `json:"data"`
		}{}

		if err := json.Unmarshal(msg, &msgStruct); err != nil {
			return err
		}

		err := c.repo.AddBalanceChange(msgStruct.Data.AccountID, msgStruct.Data.BalanceChanging, msgStruct.Data.Timestamp)
		if err != nil {
			return err
		}

		return nil
	}

	go c.consume("analytics.balance.changed", "analytics/balance.changed", handler)
}
