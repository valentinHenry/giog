package io

import (
	"context"
	"github.com/valentinHenry/giog/tuples"
	v "github.com/valentinHenry/giog/void"
)

func runEffect[A any](universe *Universe, io IO[A]) (ioRes *A, ioError Cause) {
	curr := io

	for curr != nil && ioError == nil {
		select {
		case <-universe.Context.Done(): // Error in context, not carrying-on
			ioError = &Cancellation{}
		default:
			// TODO find a way not to make a useless function call
			// Using type(~pattern~) matching with any instead of run is impossible
			res := curr.run(universe)
			curr = res.curr
			ioError = res.ioError
			ioRes = res.ioRes
		}
	}

	return ioRes, ioError
}

func (s *_IOSuccess[A, B]) run(universe *Universe) runResult[B] {
	res, prevErr := runEffect(universe, s.previous)

	if prevErr != nil {
		return retErr[B](prevErr.appendTraceIfNecessary(s.trace))
	}
	return retCurr(s.onSuccess(*res))
}
func (s *_IOSuccess[A, B]) UnsafeRun() (B, error) {
	return unsafeRun[B](s)
}

func (f *_IOFailure[A]) run(universe *Universe) runResult[A] {
	res, prevErr := runEffect(universe, f.previous)

	if prevErr == nil {
		return retRes(res)
	}

	switch prevErr.(type) {
	case *Cancellation:
		return retErr[A](prevErr)
	default:
		return retCurr[A](f.onFailure(prevErr)) // append ?
	}
}
func (f *_IOFailure[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](f)
}

func (sf *_IOSuccessFailure[A, B]) run(universe *Universe) runResult[B] {
	res, prevErr := runEffect(universe, sf.previous)

	if prevErr != nil {
		switch prevErr.(type) {
		case *Cancellation:
			return retErr[B](prevErr)
		default:
			return retCurr(sf.onFailure(prevErr))
		}
	}

	return retCurr(sf.onSuccess(*res))
}
func (f *_IOSuccessFailure[A, B]) UnsafeRun() (B, error) {
	return unsafeRun[B](f)
}

func (s *_IOExitSuccess[A]) run(_ *Universe) runResult[A] {
	res := s.success
	return retRes(&res)
}
func (f *_IOExitSuccess[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](f)
}

func (e *_IOExitError[A]) run(_ *Universe) runResult[A] {
	return retErr[A](e.cause)
}
func (f *_IOExitError[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](f)
}

func (s *_IOSync[A]) run(universe *Universe) runResult[A] {
	res := s.eval(universe.Context)
	return retRes(&res)
}
func (f *_IOSync[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](f)
}

func (us *_IOUniverseSwitch[A]) run(universe *Universe) runResult[A] {
	tmpUniverse := us.get(universe)
	ioRes, ioError := runEffect(tmpUniverse, us.withUniverses(universe, tmpUniverse))
	us.release(tmpUniverse)
	return runResult[A]{
		curr:    nil,
		ioRes:   ioRes,
		ioError: ioError,
	}
}
func (f *_IOUniverseSwitch[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](f)
}

func (wl *_IOWhileLoop) run(universe *Universe) runResult[v.Void] {
	cRes, cErr := runEffect(universe, wl.cond)
	for ; cErr == nil && *cRes; cRes, cErr = runEffect(universe, wl.cond) {
		select {
		case <-universe.Context.Done(): // Error in context, not carrying-on
			return retErr[v.Void](makeCancellationCause())
		default:
			if _, err := runEffect(universe, wl.do); err != nil {
				return retErr[v.Void](err.appendTraceIfNecessary(wl.trace))
			}
		}
	}

	if cErr != nil {
		return retErr[v.Void](cErr.appendTraceIfNecessary(wl.trace))
	}

	return retRes(&v.Void{})
}
func (wl *_IOWhileLoop) UnsafeRun() (v.Void, error) {
	return unsafeRun[v.Void](wl)
}

func (oc *_IOOnCancel[A]) run(universe *Universe) runResult[A] {
	res, prevErr := runEffect(universe, oc.previous)

	if prevErr != nil {
		switch prevErr.(type) {
		case *Cancellation:
			return retCurr(oc.onCancel)
		default:
			return retErr[A](prevErr.appendTraceIfNecessary(oc.trace))
		}
	}

	return retRes(res)
}
func (oc *_IOOnCancel[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](oc)
}

func (i *_IOInterruptFast[A]) run(universe *Universe) runResult[A] {
	resChan := make(chan tuples.T2[*A, Cause])

	go func(resChan chan tuples.T2[*A, Cause], universe *Universe) {
		res, err := runEffect(universe, i.io)
		resChan <- tuples.Of2(res, err)
	}(resChan, universe)

	select {
	case <-universe.Context.Done(): // Error in context, not carrying-on
		return retErr[A](makeCancellationCause())
	case yielded := <-resChan:
		res, err := yielded.Values()
		if err != nil {
			return retErr[A](err.appendTraceIfNecessary(i.trace))
		}
		return retRes(res)
	}
}
func (i *_IOInterruptFast[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](i)
}

func (a *_IOAsync[A]) run(universe *Universe) runResult[v.Void] {
	asyncUniverse := universe.CloneWithContext(a.ctx)

	a.runAsync(func() (A, Cause) {
		if !a.forgettable {
			universe.waitForMe()
			defer universe.done()
		}
		a, c := runEffect(asyncUniverse, a.io)
		return *a, c
	})

	return retRes(&v.Void{})
}

func (a *_IOAsync[A]) UnsafeRun() (v.Void, error) {
	return unsafeRun[v.Void](a)
}

type runResult[A any] struct {
	curr    IO[A]
	ioRes   *A
	ioError Cause
}

func retCurr[A any](curr IO[A]) runResult[A] {
	return runResult[A]{
		curr:    curr,
		ioRes:   nil,
		ioError: nil,
	}
}

func retRes[A any](res *A) runResult[A] {
	return runResult[A]{
		curr:    nil,
		ioRes:   res,
		ioError: nil,
	}
}

func retErr[A any](err Cause) runResult[A] {
	return runResult[A]{
		curr:    nil,
		ioRes:   nil,
		ioError: err,
	}
}

func unsafeRun[A any](io IO[A]) (A, error) {
	defaultUnsafeUniverse := &Universe{
		Context:         context.Background(),
		Uninterruptible: context.Background(),
		waitForMe:       func() {},
		done:            func() {},
	}

	ret, err := runEffect(defaultUnsafeUniverse, io)
	var dummyA A
	if err != nil {
		return dummyA, err.Cause()
	}
	return *ret, nil
}
