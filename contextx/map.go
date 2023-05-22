package contextx

import "sync"

type ContextValue struct {
	V map[any]interface{}

	M sync.RWMutex
}

func newContextValue() *ContextValue {
	return &ContextValue{
		V: make(map[any]interface{}),
	}
}
