package gogoing

import (
	"net"
	"log"
	"github.com/satori/go.uuid"
)

type socketAccepter struct {
	*peer
	*sessionManager
	listener net.Listener
	running bool
}

func (self *socketAccepter) Start(address string) Peer {
	ln, err := net.Listen("tcp", address)
	self.listener = ln

	if err!=nil {
		log.Printf("#listen failed (%s) %v", self.peer.Name(), err.Error())
	}

	self.running = true

	log.Printf("#listen(%s) %s", self.peer.Name(), address)

	go func() {
		for self.running {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("#accept failed(%s) %v", self.peer.Name(), err.Error())
				break
			}

			go func() {
				ses := newSession(conn, self)

				self.sessionManager.AddSession(ses)

				ses.OnClose = func() {
					self.sessionManager.RemoveSession(ses)
				}

				log.Printf("#accepted(%s) sessionId: %d", self.peer.Name(), ses.ID())
				// 发送链接事件
				//ses.recvQueue.Post(Event{Type:CONNECT_EVENT, Sess:ses})
			}()
		}
	}()
	return self
}

func (self *socketAccepter) Stop(ch chan bool) {
	if !self.running {
		return
	}
	self.running = false
	self.listener.Close()
	ch<-self.running
}

func NewAcceptor() Peer {
	self := &socketAccepter{
		peer: newPeer(uuid.NewV4().String()),
		sessionManager: newSessionManager(),
	}
	return self
}


