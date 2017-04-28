package component

import (
	"reflect"
	"db"
	"sync"
	"time"
	"fmt"
)

type Component struct {
	mu           *sync.Mutex
	componentMap map[reflect.Type]db.Proxy
	ticker       *time.Ticker
}

func (self *Component) CreateIfNotExist(proxyType reflect.Type) db.Proxy {
	self.mu.Lock()
	defer self.mu.Unlock()
	_, ok := self.componentMap[proxyType]
	if !ok {
		proxy := reflect.New(proxyType).Interface()
		//self.componentMap[proxyType] = (db.Proxy)(unsafe.Pointer(proxy))
		self.componentMap[proxyType] = proxy.(db.Proxy)
	}
	if self.componentMap[proxyType].Timer() ==nil {
		self.componentMap[proxyType].SetTimer(time.AfterFunc(time.Minute*2, func() {
			self.Remove(proxyType)
		}))
	} else {
		self.componentMap[proxyType].Timer().Reset(time.Minute*2)
	}
	return self.componentMap[proxyType]
}

func (self *Component) Remove(proxyType reflect.Type) {
	self.mu.Lock()
	defer self.mu.Unlock()
	delete(self.componentMap, proxyType)
	fmt.Println("remove component of ", proxyType.String())
}

func (self *Component) autoSave() {
	go func() {
		for {
			<-self.ticker.C
			self.save()
			fmt.Println("ticker save per 10 second...")
		}
	}()
}

func (self *Component) Stop() {
	self.mu.Lock()
	self.ticker.Stop()
	defer self.mu.Unlock()
	for _, v := range self.componentMap {
		v.Save()
	}
}

func (self *Component) save() {
	self.mu.Lock()
	defer self.mu.Unlock()
	for _, v := range self.componentMap {
		v.Save()
	}
}

func NewComponent() (component *Component) {
	component = &Component{componentMap: make(map[reflect.Type]db.Proxy), mu: new(sync.Mutex), ticker: time.NewTicker(time.Second * 10)}
	component.autoSave()
	return
}
