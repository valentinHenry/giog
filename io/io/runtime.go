package io

import (
	"context"
	"errors"
	"sync"
)

type Universe struct {
	Context         context.Context
	Uninterruptible context.Context
	waitForMe       func()
	done            func()
}

func (u *Universe) CloneWithContext(c context.Context) *Universe {
	return &Universe{
		Context:         c,
		Uninterruptible: u.Uninterruptible,
		waitForMe:       u.waitForMe,
		done:            u.done,
	}
}

// EvalAsync execute the IO and returns the sync result A or error
// Unlike EvalSync, it does not wait on any async non-waited effects in the runtime
func EvalAsync[A any](io IO[A]) (A, error) {
	universe := &Universe{
		Context:         context.Background(),
		Uninterruptible: context.Background(),
		waitForMe:       func() {},
		done:            func() {},
	}

	v, cause := runEffect(universe, io)
	if cause != nil {
		return *v, errors.New(cause.appendTraceIfNecessary(getTrace(1)).sPrettyPrint())
	}
	return *v, nil
}

// EvalSync executes the IO and waits for all async functions to end.
func EvalSync[A any](io IO[A]) (A, error) {
	wg := sync.WaitGroup{}

	universe := &Universe{
		Context:         context.Background(),
		Uninterruptible: context.Background(),
		waitForMe:       func() { wg.Add(1) },
		done:            func() { wg.Done() },
	}

	v, cause := runEffect(universe, io)

	if cause != nil {
		return *v, errors.New(cause.appendTraceIfNecessary(getTrace(1)).sPrettyPrint())
	}

	wg.Wait()
	return *v, nil
}
