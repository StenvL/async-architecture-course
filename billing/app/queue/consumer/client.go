package consumer

import (
	"github.com/StenvL/async-architecture-course/billing/app/queue/client"
	"github.com/StenvL/async-architecture-course/billing/app/queue/producer"
	"github.com/StenvL/async-architecture-course/billing/app/repository"
)

type msgHandler func(msg []byte) error

type Client struct {
	mq       client.Client
	repo     repository.Repository
	producer producer.Client
}

func New(mq client.Client, repo repository.Repository, producer producer.Client) Client {
	return Client{
		mq:       mq,
		repo:     repo,
		producer: producer,
	}
}
