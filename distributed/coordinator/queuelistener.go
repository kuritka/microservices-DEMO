package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/streadway/amqp"
	"microservices/distributed/dto"
	"microservices/distributed/qutils"
)

const url = "amqp://guest:guest@localhost:5672"

type QueueListener struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//check if we registered each sensor to do not register it twice
	sources    map[string]<-chan amqp.Delivery
	aggregator *EventAggregator
}

func NewQueueListener(aggregator *EventAggregator) *QueueListener {
	listener := QueueListener{
		sources:    make(map[string]<-chan amqp.Delivery),
		aggregator: aggregator,
	}

	listener.conn, listener.channel = qutils.GetChannel(url)

	return &listener
}

func (listener *QueueListener) ListenForNewSource() {
	//rabbit creates queue and attach some name when it is empty
	q := qutils.GetQueue("", listener.channel, true)

	//rebinding queue to fanout model
	//+ which queue to bind?
	//+ fanout exchanges ignores this field
	//+ true = when binding doesnt succeed than close channel
	//+ non need args param
	listener.channel.QueueBind(q.Name, "", "amq.fanout", false, nil)

	msgs, _ := listener.channel.Consume(q.Name, "", true, false, false, false, nil)

	listener.DiscoverSensors()

	//for loop waiting for messages in msg channel
	for msg := range msgs {
		listener.aggregator.PublishEvent("DataSourceDiscovered", string(msg.Body))
		//msg came so new sensor is online and binds into the system
		//in order to receive those messages
		//msg.body has queue name for just activated sensor
		sourceChannel, _ := listener.channel.Consume(string(msg.Body), "", true, false, false, false, nil)

		//has new message already been registered ?
		if listener.sources[string(msg.Body)] == nil {
			listener.sources[string(msg.Body)] = sourceChannel

			go listener.AddListener(sourceChannel)
		}

	}
}

func (listener *QueueListener) AddListener(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		r := bytes.NewReader(msg.Body)
		d := gob.NewDecoder(r)
		sd := new(dto.SensorMessage)
		d.Decode(sd)

		fmt.Printf("received message : %v \n", sd)

		ed := EventData{
			Name:      sd.Name,
			Value:     sd.Value,
			Timestamp: sd.Timestamp,
			Id:        sd.Id,
		}

		listener.aggregator.PublishEvent("MessageReceived_"+msg.RoutingKey, ed)
	}
}

func (listener *QueueListener) DiscoverSensors() {
	listener.channel.ExchangeDeclare(
		qutils.SensorDiscoveryExchange,
		"fanout", //could be fanout header and something else...
		false,    //setup durable queue or not ?
		false,    //atodelete exchange when nobody is present?
		false,    //true = reject external publishing requests, advanced broker scenarios
		false,
		nil,
	)

	listener.channel.Publish(qutils.SensorDiscoveryExchange, "", false, false, amqp.Publishing{})
}
