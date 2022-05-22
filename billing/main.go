package main

import (
	"log"

	"github.com/StenvL/async-architecture-course/billing/app/api"
	"github.com/StenvL/async-architecture-course/billing/app/queue/consumer"
	"github.com/StenvL/async-architecture-course/billing/app/queue/producer"

	"github.com/StenvL/async-architecture-course/billing/app/repository"

	mqclient "github.com/StenvL/async-architecture-course/billing/app/queue/client"
	"github.com/StenvL/async-architecture-course/billing/app/queue/declare"

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
	if err = declare.New(mq).DeclareTaskQueues(); err != nil {
		log.Fatal(err)
	}

	msgProducer := producer.New(mq)

	cons := consumer.New(mq, repo, msgProducer)
	if err = cons.ConsumeUserEvents(); err != nil {
		log.Fatal(err)
	}
	if err = cons.ConsumeTaskEvents(); err != nil {
		log.Fatal(err)
	}

	server := api.New(repo, msgProducer)
	server.SetupRoutes()
	server.Start()
}
