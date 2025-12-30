package domain

import "sync"

type SecureMap[K any, V any] struct {
	sMap sync.Map
}

func (sm *SecureMap[K, V]) Set(key K, value V) {
	sm.sMap.Store(key, value)
}
func (sm *SecureMap[K, V]) Delete(key K) {
	sm.sMap.Delete(key)
}
func (sm *SecureMap[K, V]) Get(key K) (V, bool) {
	val, ok := sm.sMap.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return val.(V), true
}
