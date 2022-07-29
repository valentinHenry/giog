package array

import (
	"errors"
	t "github.com/valentinHenry/giog/utils/tuples"
)

func Map[T1, T2 any](arr []T1, mapFn func(T1) T2) []T2 {
	return _map(arr, make([]T2, len(arr)), func(_ int, obj T1) T2 { return mapFn(obj) })
}

func MapK[T1, T2 any](mapFn func(T1) T2) func([]T1) []T2 {
	return func(arr []T1) []T2 {
		return _map(arr, make([]T2, len(arr)), func(_ int, obj T1) T2 { return mapFn(obj) })
	}
}

func MapWithIndex[T1, T2 any](arr []T1, mapFn func(int, T1) T2) []T2 {
	return _map(arr, make([]T2, len(arr)), mapFn)
}

func MapWithIndexK[T1, T2 any](mapFn func(int, T1) T2) func([]T1) []T2 {
	return func(arr []T1) []T2 {
		return _map(arr, make([]T2, len(arr)), mapFn)
	}
}

func MapInPlace[T1 any](arr []T1, mapFn func(T1) T1) []T1 {
	return _map(arr, arr, func(_ int, obj T1) T1 { return mapFn(obj) })
}

func MapInPlaceWithIndex[T1 any](arr []T1, mapFn func(int, T1) T1) []T1 {
	return _map(arr, arr, mapFn)
}

func FlatMap[T1, T2 any](arr []T1, mapFn func(T1) []T2) []T2 {
	return _flatMap(arr, func(_ int, t T1) []T2 { return mapFn(t) })
}

func FlatMapK[T1, T2 any](mapFn func(T1) []T2) func([]T1) []T2 {
	return func(arr []T1) []T2 {
		return _flatMap(arr, func(_ int, t T1) []T2 { return mapFn(t) })
	}
}

func FlatMapWithIndex[T1, T2 any](arr []T1, mapFn func(int, T1) []T2) []T2 {
	return _flatMap(arr, mapFn)
}

func FlatMapWithIndexK[T1, T2 any](mapFn func(int, T1) []T2) func([]T1) []T2 {
	return func(arr []T1) []T2 {
		return _flatMap(arr, mapFn)
	}
}

func Filter[T any](arr []T, filterFn func(T) bool) []T {
	return _filter(arr, func(_ int, t T) bool { return filterFn(t) })
}

func FilterK[T any](filterFn func(T) bool) func([]T) []T {
	return func(arr []T) []T {
		return _filter(arr, func(_ int, t T) bool { return filterFn(t) })
	}
}

func FilterWithIndex[T any](arr []T, filterFn func(int, T) bool) []T {
	return _filter(arr, filterFn)
}

func FilterWithIndexK[T any](arr []T, filterFn func(int, T) bool) func([]T) []T {
	return func(arr []T) []T {
		return _filter(arr, filterFn)
	}
}

func ZipWithIndex[T any](arr []T) []t.T2[int, T] {
	dst := make([]t.T2[int, T], len(arr))
	for i, c := range arr {
		dst[i] = t.Of2(i, c)
	}
	return dst
}

func ZipWithIndexK[T any]() func([]T) []t.T2[int, T] {
	return func(arr []T) []t.T2[int, T] {
		return ZipWithIndex(arr)
	}
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

func _flatMap[T1, T2](src []T1, mapFn func(int, T1) []T2) []T2 {
	var res []T2
	for i, c := range src {
		res = append(res, mapFn(i, c)...)
	}
	return res
}

func _filter[T any](src []T, filterFn func(int, T) bool) []T {
	var dst []T
	for i, c := range src {
		if filterFn(i, c) {
			dst = append(dst, c)
		}
	}
	return dst
}
