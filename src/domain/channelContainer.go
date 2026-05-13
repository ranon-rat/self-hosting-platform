package domain

import "sync"

type SafeChan[T any] struct {
	ch   chan T
	once sync.Once
}

func NewSafeChan[T any](maxBuffer int) *SafeChan[T] {
	channel := make(chan T, maxBuffer)
	return &SafeChan[T]{
		ch: channel,
	}

}

// this allows a more secure execution, since the channels can be reused(for practical usages)
func (sc *SafeChan[T]) Close() {
	sc.once.Do(func() {
		close(sc.ch)
	})
}

func (sc *SafeChan[T]) Send(v T) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = true
		}
	}()
	sc.ch <- v
	return
}

func (sc *SafeChan[T]) Receive() (T, bool) {
	v, ok := <-sc.ch
	return v, ok
}

// it adds to the buffer, if its not consumed it ignores that.
func (sc *SafeChan[T]) AppendSend(v T) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = true
		}
	}()
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
