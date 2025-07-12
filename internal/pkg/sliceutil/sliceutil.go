package sliceutil

import "sort"

func SafeGet[T any](s []T, idx int) (T, bool) {
	var zero T
	if idx < 0 || idx >= len(s) {
		return zero, false
	}
	return s[idx], true
}

func Map[T any, R any](s []T, mapper func(T) R) []R {
	if len(s) == 0 {
		return nil
	}
	out := make([]R, len(s))
	for i, v := range s {
		out[i] = mapper(v)
	}
	return out
}

func Filter[T any](s []T, pred func(T) bool) []T {
	if len(s) == 0 {
		return nil
	}
	out := make([]T, 0, len(s))
	for _, v := range s {
		if pred(v) {
			out = append(out, v)
		}
	}
	return out
}

func Contains[T any](s []T, target T, eq func(a, b T) bool) bool {
	for _, v := range s {
		if eq(v, target) {
			return true
		}
	}
	return false
}

func Reduce[T any, R any](s []T, init R, combiner func(R, T) R) R {
	acc := init
	for _, v := range s {
		acc = combiner(acc, v)
	}
	return acc
}

func SortCopy[T any](s []T, less func(a, b T) bool) []T {
	if len(s) == 0 {
		return nil
	}
	cp := make([]T, len(s))
	copy(cp, s)
	sort.Slice(cp, func(i, j int) bool { return less(cp[i], cp[j]) })
	return cp
}

func MapIdx[T any, R any](s []T, mapper func(idx int, v T) R) []R {
	if len(s) == 0 {
		return nil
	}
	out := make([]R, len(s))
	for i, v := range s {
		out[i] = mapper(i, v)
	}
	return out
}

func TryMap[T any, R any](s []T, mapper func(v T) (R, error)) ([]R, error) {
	return TryMapIdx(s, func(_ int, v T) (R, error) { return mapper(v) })
}

func TryMapIdx[T any, R any](s []T, mapper func(idx int, v T) (R, error)) ([]R, error) {
	if len(s) == 0 {
		return nil, nil
	}
	out := make([]R, len(s))
	for i, v := range s {
		r, err := mapper(i, v)
		if err != nil {
			return nil, err
		}
		out[i] = r
	}
	return out, nil
}
