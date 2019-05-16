package controller

import (
	"bytes"
	"encoding/gob"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"microservices/distributed/dto"
	"microservices/distributed/qutils"
	"microservices/distributed/web/model"
	"net/http"
	"sync"
)

var url = "amqp://guest:guest@localhost:5672"

type websocketController struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	sockets  []*websocket.Conn
	mutex    sync.Mutex
	upgrader websocket.Upgrader
}

//generic message
type message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func newWebsocketController() *websocketController {
	wsc := new(websocketController)

	wsc.conn, wsc.ch = qutils.GetChannel(url)

	wsc.upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

	go wsc.listenForSources()

	go wsc.listenForMessages()

	return wsc
}

func (wsc *websocketController) handeMessage(w http.ResponseWriter, r *http.Request) {
	socket, _ := wsc.upgrader.Upgrade(w, r, nil)
	wsc.addSocket(socket)
	go wsc.listenForDiscoveryRequests(socket)
}

func (wsc *websocketController) addSocket(conn *websocket.Conn) {
	wsc.mutex.Lock()
	//append is not threadsafe
	wsc.sockets = append(wsc.sockets)
	wsc.mutex.Unlock()
}

func (wsc *websocketController) listenForDiscoveryRequests(socket *websocket.Conn) {
	for {
		msg := message{}
		err := socket.ReadJSON(&msg)
		if err != nil {
			wsc.removeSocket(socket)
			break
		}

		if msg.Type == "discover" {
			wsc.ch.Publish("", qutils.WebappDiscoveryQueue, false, false,
				amqp.Publishing{}) //presnce of a message is trigger (not content).
			// Sending empy object is picked up by one of coordinators will respond by list sensors
		}
	}
}

func (wsc *websocketController) removeSocket(socket *websocket.Conn) {
	wsc.mutex.Lock()
	socket.Close()

	for i := range wsc.sockets {
		if wsc.sockets[i] == socket {
			//remove socket if I find it
			wsc.sockets = append(wsc.sockets[:i], wsc.sockets[i+1:]...)
		}
	}
	wsc.mutex.Unlock()
}

//sending message from rabbit down to the client
//listen for messages comming from msg broker and rsponse to discovery req
func (wsc *websocketController) listenForSources() {
	q := qutils.GetQueue("", wsc.ch, true)
	wsc.ch.QueueBind(
		q.Name,
		"",
		qutils.WebappSourceExchange,
		false,
		nil,
	)

	msgs, _ := wsc.ch.Consume(q.Name, "", true, false, false, false, nil)

	for msg := range msgs {
		sensor, _ := model.GetSensorByName(string(msg.Body))
		wsc.sendMessage(message{Type: "source", Data: sensor})
	}
}

func (wsc *websocketController) sendMessage(msg message) {
	socketsToRemove := []*websocket.Conn{}
	for _, socket := range wsc.sockets {
		err := socket.WriteJSON(msg)
		if err != nil {
			socketsToRemove = append(socketsToRemove, socket)
		}
	}

	for _, socket := range socketsToRemove {
		wsc.removeSocket(socket)
	}
}

func (wsc *websocketController) listenForMessages() {
	q := qutils.GetQueue("", wsc.ch, true)
	wsc.ch.QueueBind(q.Name, "", qutils.WebappReadingsExchange, false, nil)

	msgs, _ := wsc.ch.Consume(q.Name, "", false, false, false, false, nil)

	for msg := range msgs {
		buf := bytes.NewBuffer(msg.Body)
		dec := gob.NewDecoder(buf)
		sm := dto.SensorMessage{}
		dec.Decode(&sm)

		wsc.sendMessage(message{Type: "reading", Data: sm})
	}

}
