package db

type Proxy interface {
	Save()
	LazyLoad(primaryKey int64)
}
