package gogoing

import "fmt"

type EventHandler interface {
	OnEvent(e *Event)

	EventType() uint8
}

type eventHandler struct {
	eventType uint8
}

func (self *eventHandler) OnEvent(e *Event) {
	fmt.Println(e.Sess.ID(),e.Type)
}

func (self *eventHandler) EventType() uint8 {
	return self.eventType
}

func NewEventHandler(eventType uint8) EventHandler {
	return &eventHandler{eventType:eventType}
}
