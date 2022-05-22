package main

import (
	"log"

	"github.com/StenvL/async-architecture-course/analytics/app/queue/consumer"

	"github.com/StenvL/async-architecture-course/analytics/app/repository"

	"github.com/StenvL/async-architecture-course/analytics/app/queue/declare"

	"github.com/StenvL/async-architecture-course/analytics/app/api"
	mqclient "github.com/StenvL/async-architecture-course/analytics/app/queue/client"

	_ "github.com/lib/pq"
)

func main() {
	repo, err := repository.New()
	if err != nil {
		log.Fatal(err)
	}

	mq, err := mqclient.New()
	if err != nil {
		log.Fatal(err)
	}

	if err = declare.New(mq).DeclareExchanges(); err != nil {
		log.Fatal(err)
	}
	if err = declare.New(mq).DeclareUserQueues(); err != nil {
		log.Fatal(err)
	}
	if err = declare.New(mq).DeclareTasksQueues(); err != nil {
		log.Fatal(err)
	}
	if err = declare.New(mq).DeclareBalanceQueues(); err != nil {
		log.Fatal(err)
	}

	if err = consumer.New(mq, repo).ConsumeUserEvents(); err != nil {
		log.Fatal(err)
	}
	if err = consumer.New(mq, repo).ConsumeTasksEvents(); err != nil {
		log.Fatal(err)
	}
	if err = consumer.New(mq, repo).ConsumeBalanceEvents(); err != nil {
		log.Fatal(err)
	}

	server := api.New(repo)
	server.SetupRoutes()
	server.Start()
}
