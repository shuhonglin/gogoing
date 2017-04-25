package util

import "errors"

type Resource interface {
	Close()
	IsClosed() bool
}

type Factory func() (Resource, error)

type acquireMessage struct {
	rch chan Resource
	ech chan error
}

type releaseMessage struct {
	r Resource
}

type closeMessage struct {

}

type ResourcePool struct {
	factory Factory
	idleResources ring
	idleCapacity int
	maxResources int
	numResources int

	acqchan chan acquireMessage
	rchan chan releaseMessage
	cchan chan closeMessage

	activeWaits []acquireMessage
}


func NewResourcePool(factory Factory, idleCapacity, maxResources int) (rp *ResourcePool) {
	rp = &ResourcePool{
		factory: factory,
		idleCapacity: idleCapacity,
		maxResources:maxResources,

		acqchan: make(chan acquireMessage),
		rchan:make(chan releaseMessage, 1),
		cchan:make(chan closeMessage, 1),
	}

	go rp.mux()
	return
}

func (rp *ResourcePool) mux() {
	loop:
	for {
		select {
		case acq := <-rp.acqchan:
			rp.acquire(acq)
		case rel := <-rp.rchan:
			if len(rp.activeWaits) !=0 {
				if !rel.r.IsClosed() { //将释放的资源分配给等待队列的第一个
					rp.activeWaits[0].rch <-rel.r
				} else { //否则忽略释放的资源并重新创建一个
					r,err:=rp.factory()
					if err != nil {
						rp.numResources--
						rp.activeWaits[0].ech<-err
					} else {
						rp.activeWaits[0].rch<-r
					}
				}
				rp.activeWaits = rp.activeWaits[1:]
			} else {
				rp.release(rel.r)
			}
		case _=<-rp.cchan:
			break loop
		}
	}
	// 清空缓存池
	for !rp.idleResources.Empty()  {
		rp.idleResources.Dequeue().Close()
	}
	for _,aw := range rp.activeWaits  {
		aw.ech<-errors.New("Resource pool closed")
	}
}

func (rp *ResourcePool) acquire(acq acquireMessage) {
	for !rp.idleResources.Empty()  { // 有空闲的资源
		r:=rp.idleResources.Dequeue()
		if !r.IsClosed() {
			acq.rch<-r
			return
		}
		rp.numResources--
	}
	if rp.maxResources !=-1 && rp.numResources>=rp.maxResources { //无空闲资源加入等待列表
		rp.activeWaits = append(rp.activeWaits, acq)
		return
	}
	r, err:=rp.factory()
	if err != nil {
		acq.ech<-err
	} else {
		rp.numResources++
		acq.rch<-r
	}
	return
}

func (rp *ResourcePool) release(resource Resource) {
	if resource!=nil || resource.IsClosed() { //直接返回
		rp.numResources--
		return
	}
	if rp.idleCapacity != -1 && rp.idleResources.Size() == rp.idleCapacity {//缓冲池已满
		resource.Close()
		rp.numResources--
		return
	}
	rp.idleResources.Enqueue(resource)
}

func (rp *ResourcePool) Acquire() (res Resource, err error) {
	acp := acquireMessage{
		rch:make(chan Resource),
		ech:make(chan error),
	}
	rp.acqchan<-acp

	select {
	case res = <-acp.rch:
	case err = <-acp.ech:
	}
	return
}

// Release() will release a resource for use by others. If the idle queue is
// full, the resource will be closed.
func (rp *ResourcePool) Release(resource Resource) {
	rel := releaseMessage{
		r: resource,
	}
	rp.rchan <- rel
}

// Close() closes all the pools resources.
func (rp *ResourcePool) Close() {
	rp.cchan <- closeMessage{}
}

// NumResources() the number of resources known at this time
func (rp *ResourcePool) NumResources() int {
	return rp.numResources
}