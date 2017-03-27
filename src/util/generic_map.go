package util

import "reflect"

type GenericMap interface {
	Get(key interface{}) interface{}

	Put(key,element interface{}) (interface{}, bool)

	Remove(key interface{}) interface{}

	Clear()

	Len() int

	Contains(key interface{}) bool

	Keys() []interface{}

	Elements() []interface{}

	ToMap() map[interface{}]interface{}

	KeyType() reflect.Type

	ElementType() reflect.Type
}
