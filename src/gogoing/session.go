package gogoing

import (
	"sync"
	"structure"
	"event"
	"fmt"
	"io"
)

type Status int
const (
	_ Status = iota
	NOT_CONNECTED
	CONNECTING
	CONNECTED
	CLOSED
)

const (
	_ uint8 = iota
	CLOSE_EVENT
	INTERNET_EVENT
	CONNECT_EVENT
	EXCEPT_CLOSE_EVENT
)

type Session interface {

	Send(*event.Event)

	Close()

	ExceptionClose()

	ID() int64

	Peer() Peer

	Dispatch(e *event.Event)
}

type session struct {
	stream EventStream  // 数据流编解码

	OnReveive func()

	OnClose func()

	id int64

	peer Peer

	endSync sync.WaitGroup

	needNotifyWrite bool

	sendList *structure.EventList

	recvQueue *structure.EventQueue

	eventDispatcher event.EventDispatcher

	status Status

}

func (self *session) ID() int64 {
	return self.id
}

func (self *session) Peer() Peer {
	return self.peer
}

func (self *session) Send(e *event.Event) {
	if e != nil {
		self.sendList.Add(e)
	}
}

func (self *session) Close() {
	// todo 发送关闭或异常的event
	self.sendList.Add(&event.Event{Type: CLOSE_EVENT, Sess:self})
}

func (self *session) ExceptionClose() {
	// todo 发送关闭或异常的event
	self.sendList.Add(&event.Event{Type: EXCEPT_CLOSE_EVENT, Sess:self})
}

func (self *session) Dispatch(e *event.Event) {
	for _, handler := range self.eventDispatcher.GetHandlers(e.Type){
		handler.OnEvent(e)
	}
}

func (self *session) sendGroutine() {

	var writeList []*event.Event
	for {
		writeList = self.sendList.PickEventList()
		willExit := false

		for _,e := range writeList {
			if e.Type == EXCEPT_CLOSE_EVENT { // 收到recvGoutine发送的异常结束事件
				willExit = true
				break
			}

			if err := self.stream.Write(e); err != nil {
				willExit = true
				break
			}
		}

		if err:= self.stream.Flush(); err != nil {
			willExit = true
		}

		if willExit {
			goto exitsendloop
		}
	}
exitsendloop:
	self.status = CLOSED
	self.needNotifyWrite = false
	self.stream.Close()
	self.endSync.Done()

}

func (self *session) recvGroutine() {
	var err error
	var e *event.Event
	for {
		e,err = self.stream.Read()
		if err != nil {
			fmt.Println("解包错误: ", e)
			break
		}
		if self.OnReveive != nil {
			self.OnReveive()
		}
		self.recvQueue.Post(e)
	}

	if self.needNotifyWrite {
		self.ExceptionClose()
	}

	self.endSync.Done()
}

func newSession(conn io.ReadWriteCloser, peer Peer, eventDispatcher event.EventDispatcher) *session {
	self := &session{
		stream : NewStream(conn),
		peer:peer,
		needNotifyWrite:true,
		sendList:structure.NewEventList(),
		recvQueue:structure.NewEventQueue(),
		eventDispatcher:eventDispatcher,
		status:NOT_CONNECTED,
	}
	self.stream.MaxPacketSize(peer.MaxPacketSize())
	self.recvQueue.StartLoop()

	self.endSync.Add(2)

	go func() {
		self.endSync.Wait()
		if self.OnClose != nil {
			self.OnClose()
		}
	}()

	go self.recvGroutine()

	go self.sendGroutine()
	return self
}
