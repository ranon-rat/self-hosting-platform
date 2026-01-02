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

// this allows a more secure execution, since the channels can be reused(for practical usages)
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

// it adds to the buffer, if its not consumed it ignores that.
func (sc *SecureChanneling[T]) AppendSend(v T) {
	select {
	case sc.ch <- v:
	default:
	}
}

// it works in the same way sync.Map ranges work.  the output must be true for it to continue
func (sc *SecureChanneling[T]) Range(f func(T) bool) {
	for out := range sc.ch {
		if !f(out) {
			return
		}
	}
}
