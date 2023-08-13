package session

import "sync"

type fetchEndPoint struct {
}

var instance *fetchEndPoint
var once sync.Once

func GetInstance() *fetchEndPoint {
	once.Do(func() {
		instance = &fetchEndPoint{}
	})
	return instance
}
