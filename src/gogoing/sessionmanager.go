package gogoing

import "sync"

type SessionManager interface {
	// 获取连接
	GetSession(sessionId int64) Session

	// 遍历连接
	VisitSession(func(Session) bool)

	// 添加session
	AddSession(ses Session)

	// 移除session
	RemoveSession(ses Session) int64

	// 连接数
	SessionCount() int
}

type sessionManager struct {
	sessionMap  map[int64]Session
	sesMapMutex *sync.RWMutex
}

func (self *sessionManager) AddSession(ses Session) {
	self.sesMapMutex.Lock()
	defer self.sesMapMutex.Unlock()

	self.sessionMap[ses.ID()] = ses
}

func (self *sessionManager) RemoveSession(ses Session) int64 {
	self.sesMapMutex.Lock()
	defer self.sesMapMutex.Unlock()

	sessionId := ses.ID()
	delete(self.sessionMap, sessionId)
	return sessionId
}

func (self *sessionManager) GetSession(sessionId int64) Session {
	self.sesMapMutex.RLock()
	defer self.sesMapMutex.RUnlock()

	v, ok := self.sessionMap[sessionId]
	if ok {
		return v
	}
	return nil
}

func (self *sessionManager) SessionCount() int {
	self.sesMapMutex.RLock()
	defer self.sesMapMutex.RUnlock()
	return len(self.sessionMap)
}

func (self *sessionManager) VisitSession(callback func(Session) bool) {
	self.sesMapMutex.RLock()
	defer self.sesMapMutex.RUnlock()
	for _, ses := range self.sessionMap {
		if !callback(ses) {
			break
		}
	}
}

func newSessionManager() *sessionManager {
	return &sessionManager{
		sessionMap:  make(map[int64]Session),
		sesMapMutex: new(sync.RWMutex),
	}
}
