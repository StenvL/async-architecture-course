package producer

import "github.com/StenvL/async-architecture-course/billing/app/queue/client"

type Client struct {
	mq client.Client
}

func New(mq client.Client) Client {
	return Client{mq: mq}
}
