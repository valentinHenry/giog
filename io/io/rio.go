package io

import (
	t "github.com/valentinHenry/giog/utils/tuples"
	v "github.com/valentinHenry/giog/utils/void"
)

type RIO[A any] interface {
	use() IO[t.T2[A, VIO]]
}

func MakeRIO[A any](acquire IO[A], release func(A) VIO) RIO[A] {
	return _MakeRIO[A](getTrace(1), acquire, release)
}

func ToRIO[A any](ioa IO[A]) RIO[A] {
	return _ToRIO(getTrace(1), ioa)
}

func PureRIO[A any](ioa A) RIO[A] {
	return _PureRIO(getTrace(1), ioa)
}

func MapRIO[A, B any](r RIO[A], mapFn func(A) B) RIO[B] {
	return _MapRIO(getTrace(1), r, mapFn)
}

func FlatMapRIO[A, B any](r RIO[A], mapFn func(A) RIO[B]) RIO[B] {
	return _FlatMapRIO(getTrace(1), r, mapFn)
}

func FlatMapRIOK[A, B any](mapFn func(A) RIO[B]) func(r RIO[A]) RIO[B] {
	_trace := getTrace(1)
	return func(r RIO[A]) RIO[B] { return _FlatMapRIO(_trace, r, mapFn) }
}

func UseRIO[A, B any](r RIO[A], fn func(A) IO[B]) IO[B] {
	return _UseRIO(getTrace(1), r, fn)
}

func _MapRIO[A, B any](_trace *Trace, r RIO[A], mapFn func(A) B) RIO[B] {
	return _FlatMapRIO(_trace, r, func(a A) RIO[B] { return _PureRIO(_trace, mapFn(a)) })
}

func _FlatMapRIO[A, B any](_trace *Trace, r RIO[A], mapFn func(A) RIO[B]) RIO[B] {
	return &_RIOBind[A, B]{
		trace:  _trace,
		source: r,
		mapFn:  mapFn,
	}
}

func _MakeRIO[A any](_trace *Trace, acquire IO[A], release func(A) VIO) RIO[A] {
	return &_RIOMake[A]{
		trace:   _trace,
		acquire: acquire,
		release: release,
	}
}

func _ToRIO[A any](_trace *Trace, ioa IO[A]) RIO[A] {
	return _MakeRIO(_trace, ioa, func(_ A) VIO { return Pure(v.Void{}) })
}

func _PureRIO[A any](_trace *Trace, ioa A) RIO[A] {
	return _ToRIO(_trace, Pure(ioa))
}

// Even though an RIO free-monad runtime would be nice to have, it is impossible (is it not?) due to the lack of
// pattern matching and generic methods or generic lambda / function declaration
func _UseRIO[A, B any](_trace *Trace, r RIO[A], fn func(A) IO[B]) IO[B] {
	return _Bracket(
		_trace,
		r.use(),
		func(tp t.T2[A, VIO]) IO[B] {
			return fn(tp.V1())
		},
		func(tp t.T2[A, VIO]) VIO {
			return tp.V2()
		},
	)
}

type _RIOMake[A any] struct {
	trace   *Trace
	acquire IO[A]
	release func(A) VIO
}

func (m *_RIOMake[A]) use() IO[t.T2[A, VIO]] {
	return Map(m.acquire, func(res A) t.T2[A, VIO] { return t.Of2(res, m.release(res)) })
}

type _RIOBind[A, B any] struct {
	trace  *Trace
	source RIO[A]
	mapFn  func(A) RIO[B]
}

func (b *_RIOBind[A, B]) use() IO[t.T2[B, VIO]] {
	_trace := b.trace // FIXME Which trace should be used here ? b.trace or one passed as parameter of use()

	return FlatMap(
		b.source.use(),
		func(tp t.T2[A, VIO]) IO[t.T2[B, VIO]] {
			a, releaseA := tp.Values()

			return _FoldIOCause(
				_trace,
				b.mapFn(a).use(),
				func(c Cause) IO[t.T2[B, VIO]] {

					if _, ok := c.Cause().(exitFastRio); ok {
						return exitError[t.T2[B, VIO]](c)
					}

					return _AndThen2(_trace, releaseA, _Raise[t.T2[B, VIO]](_trace, exitFastRio{&c}))
				},
				func(tp t.T2[B, VIO]) IO[t.T2[B, VIO]] {
					b, releastB := tp.Values()
					return _Pure(t.Of2(b, AndThen2(releastB, releaseA)))
				},
			)
		},
	)
}

type exitFastRio struct {
	reason *Cause
}

func (e exitFastRio) Error() string { return "exit rio" }
