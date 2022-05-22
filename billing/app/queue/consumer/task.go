package consumer

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/StenvL/async-architecture-course/billing/app/queue/producer"

	"github.com/shopspring/decimal"

	"github.com/google/uuid"

	"github.com/StenvL/async-architecture-course/billing/app/model"
)

func (c Client) ConsumeTaskEvents() error {
	c.taskCreatedConsumer()
	c.taskCompletedConsumer()
	c.tasksShuffledConsumer()

	return nil
}

func (c Client) taskCreatedConsumer() {
	handler := func(msg []byte) error {
		type eventPayload struct {
			ID          uuid.UUID `json:"public_id"`
			Title       string    `json:"title"`
			Key         string    `json:"key"`
			Status      string    `json:"status"`
			Created     time.Time `json:"created"`
			Description string    `json:"description"`
			Assignee    int       `json:"assignee"`
		}
		msgStruct := struct {
			Version int          `json:"event_version"`
			Data    eventPayload `json:"data"`
		}{}

		if err := json.Unmarshal(msg, &msgStruct); err != nil {
			return err
		}

		if msgStruct.Version != 2 {
			return nil
		}

		cost := decimal.New(int64(rand.Intn(100)), 0)
		reward := cost.Mul(decimal.New(3, 0))
		accID, err := c.repo.CreateTask(model.Task{
			ID:          msgStruct.Data.ID,
			Title:       msgStruct.Data.Title,
			Key:         msgStruct.Data.Key,
			Status:      msgStruct.Data.Status,
			Created:     msgStruct.Data.Created,
			Description: msgStruct.Data.Description,
			Assignee:    msgStruct.Data.Assignee,
			Cost:        cost,
			Reward:      reward,
		})
		if err != nil {
			return err
		}

		if err = c.producer.BalanceChanged(producer.BalanceChanged{
			AccountID:       accID,
			BalanceChanging: cost.Neg(),
			Timestamp:       time.Now(),
		}); err != nil {
			return err
		}

		if err = c.producer.TaskEstimated(producer.TaskEstimated{
			ID:        msgStruct.Data.ID,
			Reward:    reward,
			Timestamp: msgStruct.Data.Created,
		}); err != nil {
			return err
		}

		return nil
	}

	go c.consume("billing.tasks.created", "billing/tasks.created", handler)
}

func (c Client) taskCompletedConsumer() {
	handler := func(msg []byte) error {
		type eventPayload struct {
			ID       uuid.UUID `json:"id"`
			Assignee int       `json:"assignee"`
		}
		msgStruct := struct {
			Data eventPayload `json:"data"`
		}{}

		if err := json.Unmarshal(msg, &msgStruct); err != nil {
			return err
		}

		accID, reward, err := c.repo.CompleteTask(msgStruct.Data.ID, msgStruct.Data.Assignee)
		if err != nil {
			return err
		}

		if err = c.producer.BalanceChanged(producer.BalanceChanged{
			AccountID:       accID,
			BalanceChanging: reward,
			Timestamp:       time.Now(),
		}); err != nil {
			return err
		}

		return nil
	}

	go c.consume("billing.tasks.completed", "billing/tasks.completed", handler)
}

func (c Client) tasksShuffledConsumer() {
	handler := func(msg []byte) error {
		type eventPayload map[uuid.UUID]int
		msgStruct := struct {
			Data eventPayload `json:"data"`
		}{}

		if err := json.Unmarshal(msg, &msgStruct); err != nil {
			return err
		}

		if err := c.repo.UpdateShuffledTasks(msgStruct.Data); err != nil {
			return err
		}

		for taskID, userID := range msgStruct.Data {
			taskCost, err := c.repo.GetTaskCost(taskID)
			if err != nil {
				return err
			}

			accountID, err := c.repo.GetAccountID(userID)
			if err != nil {
				return err
			}

			if err = c.producer.BalanceChanged(producer.BalanceChanged{
				AccountID:       accountID,
				BalanceChanging: taskCost.Neg(),
				Timestamp:       time.Now(),
			}); err != nil {
				return err
			}
		}

		return nil
	}

	go c.consume("billing.tasks.shuffled", "billing/tasks.shuffled", handler)
}
