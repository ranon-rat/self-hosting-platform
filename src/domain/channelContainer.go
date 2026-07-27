package domain

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type SafeChan[T any] struct {
	ch     chan T
	once   sync.Once
	closed atomic.Bool
}

func NewSafeChan[T any](maxBuffer int) *SafeChan[T] {
	channel := make(chan T, maxBuffer)
	return &SafeChan[T]{
		ch: channel,
	}

}
func (sc *SafeChan[T]) Chan() chan T {
	return sc.ch
}

// this allows a more secure execution, since the channels can be reused(for practical usages)
func (sc *SafeChan[T]) Close() {
	sc.once.Do(func() {
		close(sc.ch)
		sc.closed.Store(true)
	})
}

func (sc *SafeChan[T]) IsClosed() bool {
	return sc.closed.Load()
}
func (sc *SafeChan[T]) Send(v T) (ok bool) {
	ok = true
	defer func() {
		if recover() != nil {
			ok = false
			sc.closed.Store(true)
		}
	}()
	if sc.closed.Load() {
		return false
	}
	fmt.Println("value being send")
	sc.ch <- v
	return
}

func (sc *SafeChan[T]) Receive() (T, bool) {
	v, ok := <-sc.ch
	return v, ok
}

// it adds to the buffer, if its not consumed it ignores that.
func (sc *SafeChan[T]) AppendSend(v T) (ok bool) {
	ok = true
	defer func() {
		if recover() != nil {
			ok = false
			sc.closed.Store(true)
		}
	}()
	if sc.closed.Load() {
		return false
	}
	select {
	case sc.ch <- v:
	default:
	}
	return
}

// it works in the same way sync.Map ranges work.  the output must be true for it to continue
func (sc *SafeChan[T]) Range(f func(T) bool) {
	for out := range sc.ch {
		if !f(out) {
			return
		}
	}
}
