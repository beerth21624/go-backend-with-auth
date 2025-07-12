package maybe

import "errors"

var ErrUnwrapEmpty = errors.New("maybe: unwrap of empty value")

type Maybe[T any] struct {
	value T
	ok    bool
}

func Just[T any](v T) Maybe[T] {
	return Maybe[T]{value: v, ok: true}
}

func Nothing[T any]() Maybe[T] {
	var zero T
	return Maybe[T]{value: zero, ok: false}
}

func (m Maybe[T]) IsPresent() bool { return m.ok }

func (m Maybe[T]) IsEmpty() bool { return !m.ok }

func (m Maybe[T]) TryGet() (T, bool) {
	if !m.ok {
		return m.value, false
	}
	return m.value, true
}

func (m Maybe[T]) UnsafeGet() T {
	if !m.ok {
		panic(ErrUnwrapEmpty)
	}
	return m.value
}

func (m Maybe[T]) OrZero() T {
	if !m.ok {
		return m.value
	}
	return m.value
}

func (m Maybe[T]) OrElse(fallback T) T {
	if m.ok {
		return m.value
	}
	return fallback
}
