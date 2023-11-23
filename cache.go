package gcache

type Cacher interface {
	Get(key string, ptrValue interface{}) (err error)
	Set(key string, value interface{}, expiration time.Duration) (err error)
	Del(key ...string) (affected int64, err error)
	Flush() (err error)
	Close() (err error)
}
