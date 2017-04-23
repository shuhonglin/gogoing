package gogoing

import (
	"time"
	"fmt"
)

type EventQueue struct {
	queue chan Event
}

func (self *EventQueue) Post(e Event) {
	fmt.Println("receive post event from session ", e.GetSess().ID())
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
			fmt.Println("process event from session ", v.GetSess().ID(), )
			v.GetSess().Dispatch(v)
		}
	}()
}

func NewEventQueue() *EventQueue {
	self := &EventQueue{
		queue: make(chan Event, 10),
	}
	return self
}
