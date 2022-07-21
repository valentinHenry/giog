package io

import (
	"context"
	dll "github.com/valentinHenry/giog/utils/collections/double_link_list"
	"sync"
)

// Deferred is an interface representing a value which may not be available yet.
//
// A deferred value can be retrieved using the Get function.
// This value can only be set once using Complete.
type Deferred[A any] interface {
	// Get returns the value contained in the Deferred.
	// It blocks semantically in case the deferred is not completed.
	Get() IO[A]

	// Complete sets the value of the Deferred.
	// It returns true if the Deferred was not completed, false otherwise.
	Complete(a A) IO[bool]
}

// MakeDeferred creates a Deferred[A] which can be used later on.
func MakeDeferred[A any]() IO[Deferred[A]] {
	return Delay(func() Deferred[A] {
		var v *A

		return &deferred[A]{
			m:       sync.Mutex{},
			v:       v,
			waiters: dll.New[chan A](),
		}
	})
}

type deferred[A any] struct {
	m       sync.Mutex
	v       *A
	waiters *dll.List[chan A]
}

func (d *deferred[A]) Get() IO[A] {
	return WithContext(func(ctx context.Context) IO[A] {
		d.m.Lock()
		if d.v != nil {
			d.m.Unlock()
			return Pure(*d.v)
		}

		waiterChan := make(chan A)
		waiter := d.waiters.PushBack(waiterChan)
		d.m.Unlock()

		select {
		case <-ctx.Done():
			d.m.Lock()
			defer d.m.Unlock()
			if d.v == nil {
				d.waiters.Remove(waiter)
			}
			return exitError[A](makeCancellationCause())
		case v := <-waiterChan:
			return Pure(v)
		}
	})
}

func (d *deferred[A]) Complete(a A) IO[bool] {
	return Delay(func() bool {
		d.m.Lock()
		if d.v != nil {
			d.m.Unlock()
			return false
		}

		d.v = &a
		d.m.Unlock()

		curr := d.waiters.Front()
		for curr != nil {
			d.waiters.Remove(curr)
			curr.Value <- a
			curr = d.waiters.Front()
		}

		return true
	})
}
