package producer

import amqp "github.com/rabbitmq/amqp091-go"

func (c Client) TaskCreated(data string) error {
	return c.produce("tasks", "created", data)
}

func (c Client) TasksShuffled(data string) error {
	return c.produce("tasks", "shuffled", data)
}

func (c Client) TaskCompleted(data string) error {
	return c.produce("tasks", "completed", data)
}

func (c Client) produce(exchange, key, data string) error {
	return c.mq.Channel.Publish(
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		},
	)
}
