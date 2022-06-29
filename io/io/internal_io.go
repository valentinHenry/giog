package io

import (
	"context"
	p "github.com/valentinHenry/giog/pipes"
	v "github.com/valentinHenry/giog/void"
)

type _IOExitError[A any] struct {
	cause Cause
}

func (e *_IOExitError[A]) Void() VIO {
	return &_IOExitError[v.Void]{cause: e.cause}
}

type _IOExitSuccess[A any] struct {
	success A
}

func (e *_IOExitSuccess[A]) Void() VIO {
	return &_IOExitSuccess[v.Void]{v.Void{}}
}

type _IOSync[A any] struct {
	trace *Trace
	eval  func(context.Context) A
}

func (e *_IOSync[A]) Void() VIO {
	return &_IOSync[v.Void]{
		trace: e.trace,
		eval:  func(ctx context.Context) v.Void { e.eval(ctx); return v.Void{} },
	}
}

type _IOUniverseSwitch[A any] struct {
	trace         *Trace
	get           func(universe *Universe) *Universe
	withUniverses func(old *Universe, new *Universe) IO[A]
	release       func(universe *Universe)
}

func (us *_IOUniverseSwitch[A]) Void() VIO {
	return _As[A, v.Void](us.trace, us, v.Void{})
}

type _IOSuccess[A, B any] struct {
	trace     *Trace
	previous  IO[A]
	onSuccess func(A) IO[B]
}

func (s *_IOSuccess[A, B]) Void() VIO {
	return _As[B, v.Void](s.trace, s, v.Void{})
}

func _IOSuccessK[A, B any](
	trace *Trace,
	onSuccess func(A) IO[B],
) func(IO[A]) IO[B] {
	return func(previous IO[A]) IO[B] {
		return &_IOSuccess[A, B]{trace, previous, onSuccess}
	}
}

type _IOFailure[A any] struct {
	trace     *Trace
	previous  IO[A]
	onFailure func(Cause) IO[A]
}

func (s *_IOFailure[A]) Void() VIO {
	return _As[A, v.Void](s.trace, s, v.Void{})
}

type _IOSuccessFailure[A1, A2 any] struct {
	trace     *Trace
	previous  IO[A1]
	onSuccess func(A1) IO[A2]
	onFailure func(Cause) IO[A2]
}

func (s *_IOSuccessFailure[A1, A2]) Void() VIO {
	return _As[A2, v.Void](s.trace, s, v.Void{})
}

type _IOWhileLoop struct {
	trace *Trace
	cond  IO[bool]
	do    VIO
}

func (wl *_IOWhileLoop) Void() VIO {
	return wl
}

type _IOOnCancel[A any] struct {
	trace    *Trace
	previous IO[A]
	onCancel IO[A]
}

func (oc *_IOOnCancel[A]) Void() VIO {
	return _As[A, v.Void](oc.trace, oc, v.Void{})
}

type _IOInterruptFast[A any] struct {
	trace *Trace
	io    IO[A]
}

func (i *_IOInterruptFast[A]) Void() VIO {
	return _As[A, v.Void](i.trace, i, v.Void{})
}

type _IOAsync[A any] struct {
	trace       *Trace
	ctx         context.Context
	io          IO[A]
	runAsync    func(func() (A, Cause))
	forgettable bool
}

func (a *_IOAsync[A]) Void() VIO {
	return a
}

func succeedSyncK[A any](_trace *Trace) func(func() A) IO[A] {
	return func(a func() A) IO[A] {
		return &_IOSync[A]{
			trace: _trace,
			eval:  func(_ context.Context) A { return a() },
		}
	}
}

func raiseK[A any](_trace *Trace) func(error) IO[A] {
	return p.Pipe2K(
		func(err error) Cause { return makeRaisedCause(_trace, &err) },
		exitError[A],
	)
}

func exitError[T any](err Cause) IO[T] {
	return &_IOExitError[T]{
		cause: err,
	}
}

func exitSuccess[T any](a T) IO[T] {
	return &_IOExitSuccess[T]{
		success: a,
	}
}
