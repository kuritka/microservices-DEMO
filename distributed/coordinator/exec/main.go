package main

import (
	"fmt"
	"microservices/distributed/coordinator"
)

var consumer *coordinator.DatabaseConsumer

func main() {
	fmt.Printf("listening....")
	aggregator := coordinator.NewEventAggreagtor()
	listener := coordinator.NewQueueListener(aggregator)
	consumer := coordinator.NewDatabaseConsumer(aggregator)

	consumer.SubscribeToDataEvent("blah")

	go listener.ListenForNewSource()

	var a string
	fmt.Scanln(&a)
}
