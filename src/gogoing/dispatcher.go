package gogoing

type EventDispatcher interface {
	AddHandler(eventType uint8, handler EventHandler)

	GetHandlers(eventType uint8) []EventHandler

	RemoveEventHandlers(eventType uint8)

	Exists(eventType uint8) bool

	Clear()

	Count() int

	CountByEventType(eventType uint8) int

}

type eventDispatcher struct {
	handlerByPeer map[uint8][]EventHandler
}

func (self *eventDispatcher) AddHandler(eventType uint8, handler EventHandler) {
	handlers, ok := self.handlerByPeer[eventType]

	if !ok {
		handlers = make([]EventHandler, 0)
	}
	handlers = append(handlers, handler)

	self.handlerByPeer[eventType] = handlers
}

func (self *eventDispatcher) GetHandlers(eventType uint8) []EventHandler {
	return self.handlerByPeer[eventType]
}

func (self *eventDispatcher) RemoveEventHandlers(eventType uint8) {
	delete(self.handlerByPeer, eventType)
}

func (self *eventDispatcher) Clear() {
	self.handlerByPeer = make(map[uint8][]EventHandler)
}

func (self *eventDispatcher) Exists(eventType uint8) bool {
	_, ok := self.handlerByPeer[eventType]
	return ok
}

func (self *eventDispatcher) Count() int {
	return len(self.handlerByPeer)
}

func (self *eventDispatcher) CountByEventType(eventType uint8) int {
	return len(self.handlerByPeer[eventType])
}

func NewEventDispatcher() EventDispatcher {
	self := &eventDispatcher{
		handlerByPeer: make(map[uint8][]EventHandler),
	}
	return self
}
