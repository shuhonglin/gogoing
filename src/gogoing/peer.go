package gogoing

import (
	"github.com/satori/go.uuid"
	"event"
)

// present a server node
// 表示一个服务器节点
type Peer interface {

	PeerID() uuid.UUID

	// 开启
	Start(address string) *Peer

	// 关闭
	Stop()

	// 名字
	SetName(name string)
	Name() string

	// Session最大包大小，超过这个将视为错误，断开链接
	SetMaxPacketSize(size int)
	MaxPacketSize() int

	// 事件分发器
	event.EventDispatcher

	// Session管理器
	SessionManager

}
