package domain

import "sync"

type SecureMap[K any, V any] struct {
	SMap sync.Map
}

func NewSecureMap[K comparable, V any]() *SecureMap[K, V] {
	return &SecureMap[K, V]{}
}

func (sm *SecureMap[K, V]) Set(key K, value V) {
	sm.SMap.Store(key, value)
}
func (sm *SecureMap[K, V]) Delete(key K) {
	sm.SMap.Delete(key)
}
func (sm *SecureMap[K, V]) Get(key K) (V, bool) {
	val, ok := sm.SMap.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return val.(V), true
}

func (sm *SecureMap[K, V]) Range(f func(K, V) bool) {
	sm.SMap.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}
