package enum

import (
	"reflect"
	"sync"
)

var (
	registry   = make(map[string][]string)
	registryMu sync.RWMutex
)

func Register[T any](labels ...string) {
	t := reflect.TypeOf((*T)(nil)).Elem().Name()
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[t] = labels
}

func GetLabels[T any]() []string {
	t := reflect.TypeOf((*T)(nil)).Elem().Name()
	registryMu.RLock()
	defer registryMu.RUnlock()
	return registry[t]
}
