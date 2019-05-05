package coordinator

import "time"

type EventRaiser interface {
	AddListener(eventName string, callback func(interface{}))
}

type EventAggregator struct {
	listeners map[string][]func(interface{})
}

type EventData struct {
	Id int
	Name string
	Value float64
	Timestamp time.Time
}


func NewEventAggreagtor() *EventAggregator {
	ea := EventAggregator{
		listeners: make(map[string][]func(interface{})),
	}
	return &ea
}


func (ea *EventAggregator) AddListener(eventName string, callback func(interface{})){
	ea.listeners[eventName] =  append(ea.listeners[eventName], callback)
}

func (ea *EventAggregator) PublishEvent(eventName string, eventData interface{}) {
	if ea.listeners[eventName] != nil {
		for _, callback := range ea.listeners[eventName] {
			callback(eventData)
		}
	}
}