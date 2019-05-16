package qutils

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//well known in entire application
const SensorDiscoveryExchange = "sensor-discovery-exchange"
const PersistReadingsQueue= "persist-readings"

//used by WebApp to let coordinator know that one of them will get a list of all of the available sources
const WebappDiscoveryQueue = "WebappDiscoveryQueue"
//names of sensors that have been discovered
const WebappSourceExchange = "WebappSourceExchange"
//sendout readings to web application
const WebappReadingsExchange = "WebappReadingsExchange"

func GetChannel(url string)(*amqp.Connection, *amqp.Channel){
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed connect to RabitMQ")
	channel, err := conn.Channel()
	failOnError(err,"Failed to open a channel")
	return conn, channel
}


//autodelete would be true for all discovery queues
func GetQueue(name string, channel *amqp.Channel, autoDelete bool) (*amqp.Queue){
	q, err :=channel.QueueDeclare( //automatically creates queue if doesnt exists
		name,					//queue name
		false,			//determines if the message should be saved to disk, messages will survive servere restart
		autoDelete ,			//what to do with messages if they dont have any active consumer, true = message wil be deleted from the queue, false = keep it
		false,			//exclusive - if same queue exists in the different channel it will fail  (false), all connections will share it (true)
		false,			//only return preexisting queue which matches providing configuration, (true) server receive an error when not find, (false) create new queue if doesnt exists on the server
		nil)
	failOnError(err,"Failed to declare queue")
	return &q
}



func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err )
		panic(fmt.Sprintf("%s: %s", msg, err ))
	}
}