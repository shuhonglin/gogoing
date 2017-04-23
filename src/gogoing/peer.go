package gogoing

import (
	"github.com/satori/go.uuid"
)

// present a server node
// 表示一个服务器节点
type Peer interface {

	PeerID() uuid.UUID

	// 开启
	Start(address string) Peer

	// 关闭
	Stop(ch chan bool)

	// 名字
	SetName(name string)
	Name() string

	// Session最大包大小，超过这个将视为错误，断开链接
	SetMaxPacketSize(maxPacketSize int)
	MaxPacketSize() int

	// 事件分发器
	EventDispatcher() EventDispatcher

	// Session管理器
	//SessionManager() SessionManager

}

type peer struct {
	peerId uuid.UUID
	name string
	maxPacketSize int
	eventDispatcher EventDispatcher
}

func (self *peer) PeerID() uuid.UUID {
	return self.peerId
}

func (self *peer) Start(address string) *Peer {
	return nil
}

func (self *peer) Stop() {
	return
}

func (self *peer) SetName(name string) {
	self.name = name
}

func (self *peer) Name() string {
	return self.name
}

func (self *peer) SetMaxPacketSize(maxPacketSize int) {
	self.maxPacketSize = maxPacketSize
}

func (self *peer) MaxPacketSize() int {
	return self.maxPacketSize
}

func (self *peer) EventDispatcher() EventDispatcher {
	return self.eventDispatcher
}

func newPeer(name string) *peer {
	return &peer{peerId:uuid.NewV4(),name:name, eventDispatcher:NewEventDispatcher()}
}
