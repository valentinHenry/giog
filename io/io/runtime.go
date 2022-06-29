package io

import (
	"context"
	"errors"
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

// InterpretAsync execute the IO and returns the sync result A or error
// Unlike InterpretSync, it does not wait on any async non-waited effects in the runtime
func InterpretAsync[A any](io IO[A]) (A, error) {
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

// InterpretSync executes the IO and waits for all async functions to end.
func InterpretSync[A any](io IO[A]) (A, error) {
	nbWaitingChan := make(chan int)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	universe := &Universe{
		Context:         ctx,
		Uninterruptible: ctx,
		waitForMe:       func() { nbWaitingChan <- 1 },
		done:            func() { nbWaitingChan <- -1 },
	}

	v, cause := runEffect(universe, io)

	if cause != nil {
		return *v, errors.New(cause.appendTraceIfNecessary(getTrace(1)).sPrettyPrint())
	}

	waitingNb := 0
	for true {
		select {
		case action := <-nbWaitingChan:
			waitingNb += action

			if waitingNb == 0 {
				return *v, nil
			}
		}
	}
	return *v, nil
}
