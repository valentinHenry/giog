package io

import (
	"context"
	"github.com/valentinHenry/giog/utils/tuples"
	v "github.com/valentinHenry/giog/utils/void"
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
			curr, ioRes, ioError = curr.run(universe)
		}
	}

	return ioRes, ioError
}

func (s *_IOSuccess[A, B]) run(universe *Universe) (curr IO[B], ioRes *B, ioError Cause) {
	res, prevErr := runEffect(universe, s.previous)

	if prevErr != nil {
		return nil, nil, prevErr.appendTraceIfNecessary(s.trace)
	}
	return s.onSuccess(*res), nil, nil
}

func (s *_IOSuccess[A, B]) UnsafeRun() (B, error) {
	return unsafeRun[B](s)
}

func (f *_IOFailure[A]) run(universe *Universe) (IO[A], *A, Cause) {
	res, prevErr := runEffect(universe, f.previous)

	if prevErr == nil {
		return nil, res, nil
	}

	switch prevErr.(type) {
	case *Cancellation:
		return nil, nil, prevErr
	default:
		return f.onFailure(prevErr), nil, nil
	}
}
func (f *_IOFailure[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](f)
}

func (sf *_IOSuccessFailure[A, B]) run(universe *Universe) (IO[B], *B, Cause) {
	res, prevErr := runEffect(universe, sf.previous)

	if prevErr != nil {
		switch prevErr.(type) {
		case *Cancellation:
			return nil, nil, prevErr
		default:
			return sf.onFailure(prevErr), nil, nil
		}
	}

	return sf.onSuccess(*res), nil, nil
}
func (sf *_IOSuccessFailure[A, B]) UnsafeRun() (B, error) {
	return unsafeRun[B](sf)
}

func (s *_IOExitSuccess[A]) run(_ *Universe) (IO[A], *A, Cause) {
	res := s.success
	return nil, &res, nil
}
func (s *_IOExitSuccess[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](s)
}

func (e *_IOExitError[A]) run(_ *Universe) (IO[A], *A, Cause) {
	return nil, nil, e.cause
}
func (e *_IOExitError[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](e)
}

func (s *_IOSync[A]) run(universe *Universe) (IO[A], *A, Cause) {
	res := s.eval(universe.Context)
	return nil, &res, nil
}
func (s *_IOSync[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](s)
}

func (us *_IOUniverseSwitch[A]) run(universe *Universe) (IO[A], *A, Cause) {
	tmpUniverse := us.get(universe)
	ioRes, ioError := runEffect(tmpUniverse, us.withUniverses(universe, tmpUniverse))
	us.release(tmpUniverse)
	return nil, ioRes, ioError
}
func (us *_IOUniverseSwitch[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](us)
}

func (wl *_IOWhileLoop) run(universe *Universe) (VIO, *v.Void, Cause) {
	cRes, cErr := runEffect(universe, wl.cond)
	for ; cErr == nil && *cRes; cRes, cErr = runEffect(universe, wl.cond) {
		select {
		case <-universe.Context.Done(): // Error in context, not carrying-on
			return nil, nil, makeCancellationCause()
		default:
			if _, err := runEffect(universe, wl.do); err != nil {
				return nil, nil, err.appendTraceIfNecessary(wl.trace)
			}
		}
	}

	if cErr != nil {
		return nil, nil, cErr.appendTraceIfNecessary(wl.trace)
	}

	return nil, &v.Void{}, nil
}
func (wl *_IOWhileLoop) UnsafeRun() (v.Void, error) {
	return unsafeRun[v.Void](wl)
}

func (oc *_IOOnCancel[A]) run(universe *Universe) (IO[A], *A, Cause) {
	res, prevErr := runEffect(universe, oc.previous)

	if prevErr != nil {
		switch prevErr.(type) {
		case *Cancellation:
			return oc.onCancel, nil, nil
		default:
			return nil, nil, prevErr.appendTraceIfNecessary(oc.trace)
		}
	}

	return nil, res, nil
}
func (oc *_IOOnCancel[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](oc)
}

func (i *_IOInterruptFast[A]) run(universe *Universe) (IO[A], *A, Cause) {
	resChan := make(chan tuples.T2[*A, Cause])

	go func(resChan chan tuples.T2[*A, Cause], universe *Universe) {
		res, err := runEffect(universe, i.io)
		resChan <- tuples.Of2(res, err)
	}(resChan, universe)

	select {
	case <-universe.Context.Done(): // Error in context, not carrying-on
		return nil, nil, makeCancellationCause()
	case yielded := <-resChan:
		res, err := yielded.Values()
		if err != nil {
			return nil, nil, err.appendTraceIfNecessary(i.trace)
		}
		return nil, res, nil
	}
}
func (i *_IOInterruptFast[A]) UnsafeRun() (A, error) {
	return unsafeRun[A](i)
}

func (a *_IOAsync[A]) run(universe *Universe) (VIO, *v.Void, Cause) {
	asyncUniverse := universe.CloneWithContext(a.ctx)

	a.runAsync(func() (A, Cause) {
		if !a.forgettable {
			universe.waitForMe()
			defer universe.done()
		}
		a, c := runEffect(asyncUniverse, a.io)
		return *a, c
	})

	return nil, &v.Void{}, nil
}

func (a *_IOAsync[A]) UnsafeRun() (v.Void, error) {
	return unsafeRun[v.Void](a)
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
