package maybe

func Map[T any, R any](m Maybe[T], mapper func(T) R) Maybe[R] {
	if m.ok {
		return Just(mapper(m.value))
	}
	return Nothing[R]()
}

func FlatMap[T any, R any](m Maybe[T], mapper func(T) Maybe[R]) Maybe[R] {
	if m.ok {
		return mapper(m.value)
	}
	return Nothing[R]()
}

func TryMap[T any, R any](m Maybe[T], mapper func(T) (R, error)) (Maybe[R], error) {
	if m.ok {
		r, err := mapper(m.value)
		if err != nil {
			return Nothing[R](), err
		}
		return Just(r), nil
	}
	return Nothing[R](), nil
}

func ToPtr[T any](m Maybe[T]) *T {
	if m.ok {
		v := m.value
		return &v
	}
	return nil
}

func FromPtr[T any](p *T) Maybe[T] {
	if p != nil {
		return Just(*p)
	}
	return Nothing[T]()
}
