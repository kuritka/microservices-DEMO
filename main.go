package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//sudo docker run -d --hostname my-rabbit --name some-rabbit -e RABBITMQ_DEFAULT_USER=guest -e RABBITMQ_DEFAULT_PASS=guest -p 15672:15672 -p 5671:5671 -p 5672:5672 -p 15671:15671 -p 25672:25672 rabbitmq:3-management

func main() {
	go client()
	go server()

	var a string
	fmt.Scanln(&a)
}

func client() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()
	//ack - we want for example when we save to db fist and later we want to ack
	//exclusive - this customer is only consumer
	msgs, err := ch.Consume(
		q.Name, //name of the queue we want to receive messages from
		"",     //uniquely identifies the connection to the queue, determine on the rabbit who is listening on the queue. Important when multiple clients receiving messages from the same queue in a direct exchange. worth in scenarios when client want to stop connection. "" server fills value automatically
		true,   //automatically acknowledge receive of the message, setted to true when receiver ppek items from queue to save resources
		true,   //this client is only consumer for this queue. throws error when other clients are already registered. Or another client tries to listen later
		false,  //prevents RabbitMQ sending messages to clients that are on the same connection as the sender
		false,  //
		nil)

	failOnError(err, "Failed register consumer")
	for msg := range msgs {
		log.Printf("received : %s ", msg.Body)
	}
}

func server() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("hello world"),
	}
	for {
		ch.Publish(
			"",     //"" = default exchange, otherwise exchange name
			q.Name, //routing key, determines in some exchange types to determine which queues should get the message
			false,  //server needs to be sure that message was delivered and when
			false,  //server needs to be sure that message was delivered and when
			msg)
	}
}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed connect to RabitMQ")
	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	q, err := channel.QueueDeclare("hello", false, false, false, false, nil)
	failOnError(err, "Failed to declare queue")
	return conn, channel, &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
