package gogoing

import (
	"fmt"
	"encoding/json"
)

type Encoder interface {
	Encode(v *Event) (data []byte, err error)
}

type Decoder interface {
	Decode(data []byte) (v *Event, err error)
}


type DefaultEncoder struct {
}

func NewDefaultEncoder() *DefaultEncoder {
	return &DefaultEncoder{}
}

func (e *DefaultEncoder) Encode(v *Event) (data []byte, err error) {
	return
}

type DefaultDecoder struct {
}

func NewDefaultDecoder() *DefaultDecoder {
	return &DefaultDecoder{}
}

func (d *DefaultDecoder) Decode(data []byte) (v *Event, err error) {
	tmpEvent := &DataEvent{}
	json.Unmarshal(data, tmpEvent)
	fmt.Println(string(tmpEvent.Data))
	v = &tmpEvent.Event
	return
}