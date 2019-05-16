package coordinator

import (
	"bytes"
	"encoding/gob"
	"github.com/streadway/amqp"
	"microservices/distributed/dto"
	"microservices/distributed/qutils"
)

type WebappConsumer struct {
	er      EventRaiser
	conn    *amqp.Connection
	ch      *amqp.Channel
	sources []string
}

func (consumer *WebappConsumer) ListenForDiscoveryRequests() {
	q := qutils.GetQueue(qutils.WebappDiscoveryQueue, consumer.ch, false)
	msgs, _ := consumer.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	for range msgs {
		for _, src := range consumer.sources {
			consumer.SendMessage(src)
		}
	}
}

func (consumer *WebappConsumer) SendMessage(src string) {
	consumer.ch.Publish(
		qutils.WebappSourceExchange,
		"",
		false,
		false,
		amqp.Publishing{Body: []byte(src)})
}

func (consumer *WebappConsumer) SubscribeDataEvent(eventName string) {
	for _, v := range consumer.sources {
		if v == eventName {
			return
		}
	}
	consumer.sources = append(consumer.sources, eventName)
	consumer.SendMessage(eventName)

	consumer.er.AddListener("Me", func(eventData interface{}) {
		ed := eventData.(EventData)
		sm := dto.SensorMessage{
			Name:      ed.Name,
			Value:     ed.Value,
			Timestamp: ed.Timestamp,
		}

		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		enc.Encode(sm)

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		//will be fanout so any listener get copy of the reading
		consumer.ch.Publish(qutils.WebappReadingsExchange, "", false, false, msg)
	})
}

func NewWebappConsumer(er EventRaiser) *WebappConsumer {
	wc := WebappConsumer{
		er: er,
	}

	wc.conn, wc.ch = qutils.GetChannel(url)
	qutils.GetQueue(qutils.PersistReadingsQueue, wc.ch, false)
	defer wc.ch.Close()

	//listening for discovery requests
	go wc.ListenForDiscoveryRequests()

	wc.er.AddListener("DataSourceDiscovered", func(eventData interface{}) {
		wc.SubscribeDataEvent(eventData.(string))
	})

	wc.ch.ExchangeDeclare(qutils.WebappSourceExchange, "fanout", false, false, false, false, nil)
	wc.ch.ExchangeDeclare(qutils.WebappReadingsExchange, "fanout", false, false, false, false, nil)

	return &wc
}
