package coordinator

import (
	"bytes"
	"encoding/gob"
	"github.com/streadway/amqp"
	"microservices/distributed/dto"
	"microservices/distributed/qutils"
	"time"
)

const maxRate = time.Second * 5

type DatabaseConsumer struct{
	raiser EventRaiser
	conn *amqp.Connection
	channel *amqp.Channel
	queue *amqp.Queue
	sources []string
}

func (consumer *DatabaseConsumer) SubscribeToDataEvent(eventName string) {
	for _,v := range consumer.sources {
		if v == eventName {
			return
		}
	}

	consumer.raiser.AddListener("MessageReceived_"+eventName, func() func(interface{}){
		prevTime := time.Unix(0,0)

		buf := new(bytes.Buffer)

		return func(eventData interface{}) {
			//casting object
			event := eventData.(EventData)
			if time.Since(prevTime) > maxRate {
				prevTime = time.Now()

				sm := dto.SensorMessage{
					Name:      event.Name,
					Value:     event.Value,
					Timestamp: event.Timestamp,
					Id:        event.Id,
				}

				buf.Reset()
				enc := gob.NewEncoder(buf)
				enc.Encode(sm)

				msg := amqp.Publishing{Body: buf.Bytes()}

				consumer.channel.Publish("",qutils.PersistReadingsQueue,false, false, msg)
			}
		}
	}())

}


func NewDatabaseConsumer(raiser EventRaiser) *DatabaseConsumer {
	consumer := DatabaseConsumer{
		raiser: raiser,
	}

	consumer.conn, consumer.channel = qutils.GetChannel(url)
	consumer.queue = qutils.GetQueue(qutils.PersistReadingsQueue,consumer.channel, false)

	consumer.raiser.AddListener("DataSourceDiscovered", func(eventData interface{}){
		consumer.SubscribeToDataEvent(eventData.(string))
	})
	return &consumer
}
