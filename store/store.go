package store


type Storer interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	Del(key string)
	Clear()
}
