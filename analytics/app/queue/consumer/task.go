package consumer

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (c Client) ConsumeTasksEvents() error {
	c.taskEstimatedConsumer()
	return nil
}

func (c Client) taskEstimatedConsumer() {
	handler := func(msg []byte) error {
		type eventPayload struct {
			ID        uuid.UUID       `json:"id"`
			Reward    decimal.Decimal `json:"reward"`
			Timestamp time.Time       `json:"timestamp"`
		}
		msgStruct := struct {
			Data eventPayload `json:"data"`
		}{}

		if err := json.Unmarshal(msg, &msgStruct); err != nil {
			return err
		}

		if err := c.repo.TaskEstimated(msgStruct.Data.ID, msgStruct.Data.Reward, msgStruct.Data.Timestamp); err != nil {
			return err
		}

		return nil
	}

	go c.consume("analytics.tasks.estimated", "analytics/tasks.estimated", handler)
}
