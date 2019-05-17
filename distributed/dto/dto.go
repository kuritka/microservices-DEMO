package dto

import (
	"time"
	//we need to encode messages for rabbit. gob is bit more efficient than jsonand pure go
	"encoding/gob"
)

type SensorMessage struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

func init() {
	gob.Register(SensorMessage{})
}
