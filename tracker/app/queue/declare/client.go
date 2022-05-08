package declare

import "github.com/StenvL/async-architecture-course/app/queue/client"

type Client struct {
	mq client.Client
}

func New(mq client.Client) Client {
	return Client{mq: mq}
}