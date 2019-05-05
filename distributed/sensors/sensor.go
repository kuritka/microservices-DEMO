package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/streadway/amqp"

	"microservices/distributed/dto"
	"microservices/distributed/qutils"
)

//go run sensor.go --help
var name = flag.String("name","sensor", "name of the sensor")
var freq = flag.Uint("freq", 5, "update frequency per sec")
var max = flag.Float64("max", .5, "max value for generated readings")
var min = flag.Float64("min", 1., "minimum value for generated readings")
var stepSize = flag.Float64("step", 0.1, "maximum allowable change per measurement")

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var value = r.Float64() * (*max - *min) + *min
var nom = (*max - *min) /2 + *min

var url = "amqp://guest:guest@localhost:5672"

func main(){
	flag.Parse()

	conn, channel := qutils.GetChannel(url)
	defer channel.Close()
	defer conn.Close() //http://localhost:15672/#/exchanges

	dataQueue := qutils.GetQueue(*name, channel, false)
	//discovery queue is known in entire system and provides new sensors
	//discoveryQueue := qutils.GetQueue(qutils.DiscoveryQueue,channel)

	publishQueueName(channel)

	//we would autoDelete queue everytime when spinup new discovery queue
	discoveryQueue := qutils.GetQueue("", channel, true)
	channel.QueueBind(discoveryQueue.Name,"", qutils.SensorDiscoveryExchange, false, nil)
	go listenForDiscoverRequest(discoveryQueue.Name, channel)


	dur,_ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	signal := time.Tick(dur)
	buf := new (bytes.Buffer)

	i := 1
	go func() {
		for range signal {
			calcValue()
			reading := dto.SensorMessage{
				Id: i,
				Name: *name,
				Value:value,
				Timestamp:time.Now(),
			}
			buf.Reset()
			enc := gob.NewEncoder(buf)  //needs to be recreated everytime when needs to be used
			enc.Encode(reading)

			msg := amqp.Publishing{
				Body: buf.Bytes(),
			}

			channel.Publish("",dataQueue.Name, false,false,msg)
			log.Printf("%s-%v %v",*name, i,value)
			i++
		}
	}()
	var value string
	fmt.Scanln(&value)
}


//if listening from amq.fanout from controller comes, tahn it push new item to fanout comming back to controllers.
//because controllers has list of existing connections, only new one will be identified
func listenForDiscoverRequest(queueName string, channel *amqp.Channel) {
	msgs,_ := channel.Consume(queueName,"",true,false,false,false,nil)
	for range msgs {
		publishQueueName(channel)
	}
}


func publishQueueName(channel *amqp.Channel){
	msg := amqp.Publishing{Body:[]byte(*name)}
	//amq.exchange - viz console -> exchange
	//key is not needed, fanout doesnt determine where the message goes
	//if
	channel.Publish("amq.fanout", "",false, false,msg)
}

func calcValue() {
	var maxStep, minStep float64

	if value <nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value ) / (*max - nom )
		minStep = -1 * *stepSize
	}
	value += r.Float64() *  (maxStep - minStep) + minStep
}