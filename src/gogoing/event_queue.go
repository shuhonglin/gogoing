package gogoing

import (
	"time"
)

type EventQueue struct {
	queue chan Event
}

func (self *EventQueue) Post(e Event) {
	self.queue <- e
}

func (self *EventQueue) DelayPost(e Event, dur time.Duration) {
	go func() {
		time.AfterFunc(dur, func() {
			self.Post(e)
		})
	}()
}

func (self *EventQueue) StartLoop() {
	go func() {
		for v:= range self.queue {
			v.Sess.Dispatch(&v)
		}
	}()
}

func NewEventQueue() *EventQueue {
	self := &EventQueue{
		queue: make(chan Event, 10),
	}
	return self
}
