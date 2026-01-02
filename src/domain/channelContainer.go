package domain

import "sync"

type SecureChanneling[T any] struct {
	ch   chan T
	once sync.Once
}

func NewSecureChanneling[T any](maxBuffer int) *SecureChanneling[T] {
	channel := make(chan T, maxBuffer)
	return &SecureChanneling[T]{
		ch: channel,
	}
}
func (sc *SecureChanneling[T]) Close() {
	sc.once.Do(func() {
		close(sc.ch)
	})
}

func (sc *SecureChanneling[T]) Send(v T) {
	sc.ch <- v
}
func (sc *SecureChanneling[T]) Receive() T {
	return <-sc.ch
}
func (sc *SecureChanneling[T]) AppendSend(v T) {
	select {
	case sc.ch <- v:
	default:
	}
}
func (sc *SecureChanneling[T]) Range(f func(T) bool) {
	for out := range sc.ch {
		if !f(out) {
			return
		}
	}
}
