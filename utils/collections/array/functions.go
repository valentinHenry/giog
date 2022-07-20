package array

import "errors"

func Map[T1, T2 any](arr []T1, mapFn func(int, T1) T2) []T2 {
	res := make([]T2, len(arr))
	return _map(arr, res, mapFn)
}

func MapInPlace[T1 any](arr []T1, mapFn func(int, T1) T1) []T1 {
	return _map(arr, arr, mapFn)
}

func _map[T1, T2 any](src []T1, dst []T2, fn func(int, T1) T2) []T2 {
	if len(src) != len(dst) {
		panic(errors.New("source and destination size are different"))
	}

	for i, c := range src {
		dst[i] = fn(i, c)
	}

	return dst
}
