package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"microservices/distributed/datamanager"
	"microservices/distributed/dto"
	"microservices/distributed/qutils"
)

var url = "amqp://guest:guest@localhost:5672"

func main(){
	fmt.Printf("datamanager listening...	")
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		qutils.PersistReadingsQueue,
		"",		//assign name for us
		false,		//true =  message automatically acknowledged and removed from the queue as soon as consumer got it. false = I'm acknowledge after DB succesfully saved record
		true,		//fails if cannot get exclusive connection to this queue - practically one db manager running in one time in the system. keeping in mind to run things in sequence
		false,
		false,
		nil)

	if err != nil {
		log.Fatalln("Failed to get access to messages %s", msgs)
	}

	for msg := range msgs {
		buf := bytes.NewReader(msg.Body)
		dec := gob.NewDecoder(buf)
		sd := &dto.SensorMessage{}
		dec.Decode(sd)


		err := datamanager.SaveReader(sd)
		if err != nil {
			log.Printf("failed to save reading fom sensor %v. Error: %s ", sd.Name, err.Error())
		} else {
			msg.Ack(false)	//acknowledge that message was processed properly and can be removed from the queue
		}
	}

}
