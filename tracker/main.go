package main

import (
	"log"

	"github.com/StenvL/async-architecture-course/tracker/app/queue/consumer"

	"github.com/StenvL/async-architecture-course/tracker/app/repository"

	"github.com/StenvL/async-architecture-course/tracker/app/queue/declare"
	"github.com/StenvL/async-architecture-course/tracker/app/queue/producer"

	"github.com/StenvL/async-architecture-course/tracker/app/api"
	mqclient "github.com/StenvL/async-architecture-course/tracker/app/queue/client"

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
	if err = declare.New(mq).DeclareTaskQueues(); err != nil {
		log.Fatal(err)
	}
	if err = consumer.New(mq, repo).ConsumeUserEvents(); err != nil {
		log.Fatal(err)
	}
	msgProducer := producer.New(mq)

	server := api.New(repo, msgProducer)
	server.SetupRoutes()
	server.Start()
}
