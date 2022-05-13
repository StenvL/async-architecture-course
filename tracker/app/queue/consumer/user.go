package consumer

import (
	"encoding/json"

	"github.com/StenvL/async-architecture-course/app/model"
	"github.com/guregu/null"
)

func (c Client) ConsumeUserEvents() error {
	c.userCreatedConsumer()
	c.userUpdatedConsumer()
	c.userDeletedConsumer()
	c.userRoleChangedConsumer()

	return nil
}

func (c Client) userCreatedConsumer() {
	handler := func(msg []byte) error {
		type eventPayload struct {
			ID    int         `json:"id"`
			Email string      `json:"email"`
			Name  null.String `json:"full_name"`
			Role  string      `json:"role"`
		}
		msgStruct := struct {
			Data eventPayload `json:"data"`
		}{}

		if err := json.Unmarshal(msg, &msgStruct); err != nil {
			return err
		}

		if err := c.repo.CreateUser(model.User{
			ID:    msgStruct.Data.ID,
			Name:  msgStruct.Data.Name,
			Email: msgStruct.Data.Email,
			Role:  msgStruct.Data.Role,
		}); err != nil {
			return err
		}

		return nil
	}

	go c.consume("users.created", "tracker/users.created", handler)
}

func (c Client) userUpdatedConsumer() {
	handler := func(msg []byte) error {
		type eventPayload struct {
			ID    int         `json:"id"`
			Email string      `json:"email"`
			Name  null.String `json:"full_name"`
		}
		msgStruct := struct {
			Data eventPayload `json:"data"`
		}{}

		if err := json.Unmarshal(msg, &msgStruct); err != nil {
			return err
		}

		if err := c.repo.UpdateUser(model.User{
			ID:    msgStruct.Data.ID,
			Name:  msgStruct.Data.Name,
			Email: msgStruct.Data.Email,
		}); err != nil {
			return err
		}

		return nil
	}

	go c.consume("users.updated", "tracker/users.updated", handler)
}
func (c Client) userDeletedConsumer() {
	handler := func(msg []byte) error {
		type eventPayload struct {
			ID int `json:"id"`
		}
		msgStruct := struct {
			Data eventPayload `json:"data"`
		}{}

		if err := json.Unmarshal(msg, &msgStruct); err != nil {
			return err
		}

		if err := c.repo.DeleteUser(msgStruct.Data.ID); err != nil {
			return err
		}

		return nil
	}

	go c.consume("users.deleted", "tracker/users.deleted", handler)
}

func (c Client) userRoleChangedConsumer() {
	handler := func(msg []byte) error {
		type eventPayload struct {
			ID   int    `json:"id"`
			Role string `json:"role"`
		}
		msgStruct := struct {
			Data eventPayload `json:"data"`
		}{}

		if err := json.Unmarshal(msg, &msgStruct); err != nil {
			return err
		}

		if err := c.repo.UpdateUserRole(msgStruct.Data.ID, msgStruct.Data.Role); err != nil {
			return err
		}

		return nil
	}

	go c.consume("users.role_changed", "tracker/users.role_changed", handler)
}

func (c Client) consume(queueName, consumerName string, handler msgHandler) error {
	msgs, _ := c.mq.Channel.Consume(
		queueName,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)

	for msg := range msgs {
		if err := handler(msg.Body); err != nil {
			_ = msg.Nack(false, false)
		}

		if err := msg.Ack(false); err != nil {
			return err
		}
	}

	return nil
}
