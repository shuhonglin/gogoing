package event

type EventHandler interface {

	OnEvent(e *Event)

	EventType() uint8

}
