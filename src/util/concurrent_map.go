package util

import (
	"reflect"
	"sync"
)

type ConcurrentMap interface {
	GenericMap
}

type GoConcurrentMap struct {
	m           map[interface{}]interface{}
	keyType     reflect.Type
	elementType reflect.Type
	mu          *sync.RWMutex
}

func NewGoConcurrentMap(keyType, elementType reflect.Type) *GoConcurrentMap {
	conMap := &GoConcurrentMap{
		m:           make(map[interface{}]interface{}),
		keyType:     keyType,
		elementType: elementType,
		mu:          new(sync.RWMutex),
	}
	return conMap
}

func (m *GoConcurrentMap) Get(key interface{}) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m[key]
}

func (m *GoConcurrentMap) CheckPair(key, element interface{}) bool {
	if key != nil && element != nil && reflect.TypeOf(key) == m.keyType && reflect.TypeOf(element) == m.elementType {
		return true
	}
	return false
}

func (m *GoConcurrentMap) Put(key, element interface{}) (interface{}, bool) {
	if !m.CheckPair(key, element) {
		return nil, false
	}
	m.mu.RLock()
	var present bool
	var oldElement interface{}
	if oldElement, present = m.m[key]; !present {
		m.mu.RUnlock()
		m.mu.Lock()
		if oldElement, present = m.m[key]; !present {
			m.m[key] = element
		}
		m.mu.Unlock()
	} else {
		m.mu.RUnlock()
	}
	return oldElement, true
}

func (m *GoConcurrentMap) Remove(key interface{}) interface{} {

	m.mu.RLock()
	var present bool
	var val interface{}
	if val, present = m.m[key]; present {
		m.mu.RUnlock()
		m.mu.Lock()
		if val, present = m.m[key]; present {
			delete(m.m, key)
		}
		m.mu.Unlock()
	} else {
		m.mu.RUnlock()
	}
	return val
}

func (m *GoConcurrentMap) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k := range m.m {
		delete(m.m, k)
	}
}

func (m *GoConcurrentMap) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.m)
}

func (m *GoConcurrentMap) Contains(key interface{}) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.m[key]
	return ok
}

func (m *GoConcurrentMap) Keys() []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]interface{}, 0, len(m.m))
	for k := range m.m {
		keys = append(keys, k)
	}
	return keys
}

func (m *GoConcurrentMap) Elements() []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	elements := make([]interface{}, 0, len(m.m))
	for _, v := range m.m {
		elements = append(elements, v)
	}
	return elements
}

func (m *GoConcurrentMap) ToMap() map[interface{}]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	newMap := make(map[interface{}]interface{})
	for k, v := range m.m {
		newMap[k] = v
	}
	return newMap
}

func (m *GoConcurrentMap) KeyType() reflect.Type {
	return m.keyType
}

func (m *GoConcurrentMap) ElementType() reflect.Type {
	return m.elementType
}
