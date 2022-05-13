package client

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

func New() (Client, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return Client{}, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return Client{}, err
	}

	return Client{
		conn:    conn,
		Channel: channel,
	}, nil
}
