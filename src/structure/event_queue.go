package structure

import (
	"event"
	"time"
)

type EventQueue struct {
	queue chan *event.Event
}

func (self *EventQueue) Post(e *event.Event) {
	self.queue <- e
}

func (self *EventQueue) DelayPost(e *event.Event, dur time.Duration) {
	go func() {
		time.AfterFunc(dur, func() {
			self.Post(e)
		})
	}()
}

func (self *EventQueue) StartLoop() {
	go func() {
		for v:= range self.queue {
			v.Sess.Dispatch(v)
		}
	}()
}

func NewEventQueue() *EventQueue {
	self := &EventQueue{
		queue: make(chan *event.Event, 10),
	}
	return self
}
