package producer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/StenvL/async-architecture-course/tracker/app/model"

	"schemaregistry"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type event struct {
	EventId      uuid.UUID   `json:"event_id"`
	EventVersion int         `json:"event_version"`
	EventName    string      `json:"event_name"`
	EventTime    string      `json:"event_time"`
	Producer     string      `json:"producer"`
	Data         interface{} `json:"data"`
}

func (c Client) TaskCreated(data model.Task) error {
	return c.produce("tasks.created", data, 1)
}

func (c Client) TasksShuffled(data map[uuid.UUID]int) error {
	return c.produce("tasks.shuffled", data, 1)
}

func (c Client) TaskCompleted(data model.TaskCompletedEvent) error {
	return c.produce("tasks.completed", data, 1)
}

func (c Client) produce(exchange string, data interface{}, version int) error {
	eventID, _ := uuid.NewUUID()

	eJSON, err := json.Marshal(event{
		EventId:      eventID,
		EventVersion: version,
		EventName:    exchange,
		EventTime:    time.Now().String(),
		Producer:     "tracker",
		Data:         data,
	})
	if err != nil {
		return fmt.Errorf("marshal event to JSON: %w", err)
	}

	valid, err := schemaregistry.Validate(exchange, eJSON, 1)
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("invalid event schema: %s", eJSON)
	}

	return c.mq.Channel.Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        eJSON,
		},
	)
}
