package db

import (
	"time"
	"db/entity"
)

type Proxy interface {
	Save()
	LazyLoad(primaryKey int64) *entity.Userinfo
	SetTimer(timer *time.Timer)
	Timer() *time.Timer
}
