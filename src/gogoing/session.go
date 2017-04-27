package gogoing

import (
	"sync"
	"fmt"
	"io"
	"component"
)

type Status int

const (
	_             Status = iota
	NOT_CONNECTED
	CONNECTING
	CONNECTED
	CLOSED
)

const (
	_                  uint8 = iota
	CLOSE_EVENT
	INTERNET_EVENT
	CONNECT_EVENT
	EXCEPT_CLOSE_EVENT
)

type Session interface {
	Send(e Event)

	Close()

	ExceptionClose()

	ID() int64

	Peer() Peer

	Dispatch(e Event)
}

type session struct {
	stream EventStream // 数据流编解码

	OnReveive func()

	OnClose func()

	id int64

	peer Peer

	endSync sync.WaitGroup

	needNotifyWrite bool

	sendList *EventList

	recvQueue *EventQueue

	components *component.Component

	status Status
}

func (self *session) ID() int64 {
	return self.id
}

func (self *session) Peer() Peer {
	return self.peer
}

func (self *session) Send(e Event) {
	if e != nil {
		self.sendList.Add(e)
	}
}

func (self *session) Close() {
	// todo 发送关闭或异常的event
	closeEvent := new(CloseEvent)
	closeEvent.Type = CLOSE_EVENT
	closeEvent.ID = 0
	closeEvent.Sess = self
	self.sendList.Add(closeEvent)
}

func (self *session) ExceptionClose() {
	// todo 发送关闭或异常的event
	//self.sendList.Add(&Event{Type: EXCEPT_CLOSE_EVENT, Sess:self})
	exceptEvent := new(ExceptEvent)
	exceptEvent.Type = EXCEPT_CLOSE_EVENT
	exceptEvent.ID = 0
	exceptEvent.Sess = self
	self.sendList.Add(exceptEvent)
}

func (self *session) Dispatch(e Event) {
	handlers := self.peer.EventDispatcher().GetHandlers(e.GetType())
	if len(handlers) > 0 {
		for _, handler := range handlers {
			handler.OnEvent(e)
		}
	} else {
		fmt.Println("unknown event type -> ", e.GetType())
	}
}

func (self *session) sendGroutine() {

	var writeList []Event
	for {
		writeList = self.sendList.PickEventList()
		willExit := false

		for _, e := range writeList {
			if e.GetType() == EXCEPT_CLOSE_EVENT { // 收到recvGoutine发送的异常结束事件
				willExit = true
				break
			}

			if err := self.stream.Write(e); err != nil {
				willExit = true
				break
			}
		}

		if err := self.stream.Flush(); err != nil {
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
	var e Event
	for {
		e, err = self.stream.Read()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端断开连接: ", err.Error())
				self.Close() // 发送客户端断开连接的event
			}
			break
		}
		if self.OnReveive != nil {
			self.OnReveive()
		}
		e.SetSess(self)
		self.recvQueue.Post(e)
	}

	if self.needNotifyWrite {
		self.ExceptionClose()
	}

	self.endSync.Done()
}

func newSession(conn io.ReadWriteCloser, peer Peer) *session {
	self := &session{
		stream:          NewStream(conn),
		peer:            peer,
		needNotifyWrite: true,
		sendList:        NewEventList(),
		recvQueue:       NewEventQueue(),
		components:      component.NewComponent(),
		status:          NOT_CONNECTED,
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
