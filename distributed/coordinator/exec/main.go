package main

import (
	"fmt"
	"microservices/distributed/coordinator"
)

var dc *coordinator.DatabaseConsumer
var wc *coordinator.WebappConsumer

func main() {
	ea := coordinator.NewEventAggregator()
	dc = coordinator.NewDatabaseConsumer(ea)
	wc = coordinator.NewWebappConsumer(ea)
	ql := coordinator.NewQueueListener(ea)
	go ql.ListenForNewSource()

	fmt.Printf("coordinator listening....")
	var a string
	fmt.Scanln(&a)
}

//package main
//
//import (
//	"fmt"
//	"microservices/distributed/coordinator"
//)
//
//var consumer *coordinator.DatabaseConsumer
//var webappConsumer *coordinator.WebappConsumer
//
//func main() {
//	fmt.Printf("listening....")
//	aggregator := coordinator.NewEventAggreagtor()
//	listener := coordinator.NewQueueListener(aggregator)
//	consumer := coordinator.NewDatabaseConsumer(aggregator)
//	webappConsumer := coordinator.NewWebappConsumer(aggregator)
//
//	consumer.SubscribeToDataEvent("blah")
//	webappConsumer.SubscribeDataEvent("blaa")
//
//	go listener.ListenForNewSource()
//
//	var a string
//	fmt.Scanln(&a)
//}
