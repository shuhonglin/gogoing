package gogoing

import (
	"fmt"
	"encoding/json"
)

type Encoder interface {
	Encode(v Event) (data []byte, err error)
}

type Decoder interface {
	Decode(data []byte) (Event, error)
}


type DefaultEncoder struct {
}

func NewDefaultEncoder() *DefaultEncoder {
	return &DefaultEncoder{}
}

func (e *DefaultEncoder) Encode(v Event) (data []byte, err error) {
	return
}

type DefaultDecoder struct {
}

func NewDefaultDecoder() *DefaultDecoder {
	return &DefaultDecoder{}
}

func (d *DefaultDecoder) Decode(data []byte) (Event, error) {
	/*if v.GetType() == CONNECT_EVENT || v.GetType()==CLOSE_EVENT || v.GetType()==EXCEPT_CLOSE_EVENT {
		fmt.Println("receive connect,close or exception event!")
		return nil, nil
	}*/
	if data==nil || len(data)==0 {
		return nil, nil
	}
	//eventType := uint8(data[0])
	/*eventId := int(data[1]<<0) | int(data[2]<<8) | int(data[3]<<16) | int(data[4]<<24)
	fmt.Println("eventId -> ", eventId)*/
	v1 := &InternetEvent{}
	json.Unmarshal(data, v1)
	fmt.Println(string(v1.Data))
	return v1, nil
}