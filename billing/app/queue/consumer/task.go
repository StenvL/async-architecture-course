package consumer

import (
	"encoding/json"
	"math/rand"
	"time"

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
			Status      string    `json:"status"`
			Created     time.Time `json:"created"`
			Description string    `json:"description"`
			Assignee    int       `json:"assignee"`
		}
		msgStruct := struct {
			Data eventPayload `json:"data"`
		}{}

		if err := json.Unmarshal(msg, &msgStruct); err != nil {
			return err
		}

		cost := decimal.New(int64(rand.Intn(100)), 0)
		task := model.Task{
			ID:          msgStruct.Data.ID,
			Title:       msgStruct.Data.Title,
			Status:      msgStruct.Data.Status,
			Created:     msgStruct.Data.Created,
			Description: msgStruct.Data.Description,
			Assignee:    msgStruct.Data.Assignee,
			Cost:        cost,
			Reward:      cost.Mul(decimal.New(3, 0)),
		}
		if err := c.repo.CreateTask(task); err != nil {
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

		if err := c.repo.CompleteTask(msgStruct.Data.ID, msgStruct.Data.Assignee); err != nil {
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

		return nil
	}

	go c.consume("billing.tasks.shuffled", "billing/tasks.shuffled", handler)
}
