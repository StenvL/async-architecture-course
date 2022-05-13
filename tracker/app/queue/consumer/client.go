package consumer

import (
	"github.com/StenvL/async-architecture-course/app/queue/client"
	"github.com/StenvL/async-architecture-course/app/repository"
)

type msgHandler func(msg []byte) error

type Client struct {
	mq   client.Client
	repo repository.Repository
}

func New(mq client.Client, repo repository.Repository) Client {
	return Client{
		mq:   mq,
		repo: repo,
	}
}
