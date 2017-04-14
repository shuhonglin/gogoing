package structure

import (
	"event"
	"sync"
)

type EventList struct {
	list  []*event.Event
	listGuard sync.Mutex

	listCond *sync.Cond
}

func (self *EventList) Add(e *event.Event) {
	self.listGuard.Lock()
	self.list = append(self.list, e)
	self.listGuard.Unlock()
	self.listCond.Signal()
}

func (self *EventList) Reset() {
	self.listGuard.Lock()
	self.list = self.list[0:0]
	self.listGuard.Unlock()
}

func (self *EventList) PickEventList() []*event.Event {
	self.listGuard.Lock()
	for len(self.list) == 0 {
		self.listCond.Wait()
	}

	dataList := make([]*event.Event, len(self.list))
	copy(dataList, self.list)
	self.listGuard.Unlock()
	return dataList
}

func NewEventList() *EventList {
	self := &EventList{}
	self.listCond = sync.NewCond(&self.listGuard)
	return self
}
