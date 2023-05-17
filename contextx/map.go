package contextx

import "sync"

type contextValue struct {
	v map[any]interface{}

	m sync.RWMutex
}

func newContextValue() *contextValue {
	return &contextValue{
		v: make(map[any]interface{}),
	}
}
