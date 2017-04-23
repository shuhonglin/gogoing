package gogoing

import "fmt"

type EventHandler interface {
	OnEvent(e Event)

	EventType() uint8
}

type eventHandler struct {
	eventType uint8
}

func (self *eventHandler) OnEvent(e Event) {
	fmt.Println(e.GetSess().ID(),e.GetType())
	if e.GetType() == INTERNET_EVENT {
		fmt.Println("on event -> ", string(e.(*InternetEvent).Data))
	}
}

func (self *eventHandler) EventType() uint8 {
	return self.eventType
}

type dataEventHandler struct {
	eventType uint8
	data []byte
}

func (self *dataEventHandler) OnEvent(e Event) {
	if e.GetType() == INTERNET_EVENT {
		fmt.Println(e.GetType())
	}
	fmt.Println(e)
}

func (self *dataEventHandler) EventType() uint8 {
	return self.eventType
}

func NewEventHandler(eventType uint8) EventHandler {
	return &eventHandler{eventType:eventType}
}
