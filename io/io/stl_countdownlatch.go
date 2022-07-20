package io

import (
	"context"
	dll "github.com/valentinHenry/giog/utils/collections/double_link_list"
	v "github.com/valentinHenry/giog/utils/void"
	r "github.com/valentinHenry/refined"
	"sync"
)

// CountDownLatch is an interface that semantically blocks any goroutines
// which waits on it. These are blocked until all defined latches are
// released.
//
// After all latches are released, the latch count is not reset
// (unlike what CyclicBarrier does with waiters). Thus Await() will not be
// blocking.
type CountDownLatch interface {
	// Release releases a latch. It does nothing in case there is no more latch.
	Release() VIO
	// Await blocks semantically until all latches are released.
	Await() VIO
}

func MakeCountDownLatch(nb r.PosInt) IO[CountDownLatch] {
	return Pure[CountDownLatch](
		&countDownLatch{
			m:                sync.Mutex{},
			remainingLatches: nb.Value(),
			waiters:          dll.List[chan any]{},
		})
}

type countDownLatch struct {
	m                sync.Mutex
	remainingLatches int
	waiters          dll.List[chan any]
}

func (cl *countDownLatch) Release() VIO {
	return Delay(func() v.Void {
		cl.m.Lock()
		defer cl.m.Unlock()

		if cl.remainingLatches == 0 {
			return v.Void{}
		}

		if cl.remainingLatches == 1 {
			curr := cl.waiters.Front()
			for curr != nil {
				close(curr.Value)
				cl.waiters.Remove(curr)
				curr = cl.waiters.Front()
			}
			cl.remainingLatches = 0
			return v.Void{}
		}

		cl.remainingLatches -= 1
		return v.Void{}
	})
}

func (cl *countDownLatch) Await() VIO {
	return WithContext(func(ctx context.Context) VIO {
		cl.m.Lock()
		if cl.remainingLatches == 0 {
			cl.m.Unlock()
			return Void()
		}

		waitingChan := make(chan any)
		waiter := cl.waiters.PushBack(waitingChan)
		cl.m.Unlock()

		select {
		case <-ctx.Done():
			cl.m.Lock()
			cl.waiters.Remove(waiter)
			cl.m.Unlock()
			return exitError[v.Void](makeCancellationCause())
		case <-waitingChan:
			return Void()
		}
	})
}
