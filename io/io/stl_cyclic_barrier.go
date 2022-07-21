package io

import (
	"context"
	dll "github.com/valentinHenry/giog/utils/collections/double_link_list"
	v "github.com/valentinHenry/giog/utils/void"
	r "github.com/valentinHenry/refined"
	"sync"
)

// A CyclicBarrier is an interface to a synchronizer which allows goroutines
// to wait for each-others at a fixed point.
type CyclicBarrier interface {
	// Await blocks semantically until enough goroutines are waiting at this barrier.
	Await() VIO
}

// MakeCyclicBarrier creates a CyclicBarrier. It takes the number of goroutines
// which must wait before all the waiters are released.
func MakeCyclicBarrier(parties r.PosInt) IO[CyclicBarrier] {
	return Pure[CyclicBarrier](
		&cyclicBarrier{
			waiters: dll.List[chan any]{},
			parties: parties,
			m:       sync.Mutex{},
		},
	)
}

type cyclicBarrier struct {
	waiters dll.List[chan any]
	parties r.PosInt
	m       sync.Mutex
}

func (cb *cyclicBarrier) Await() VIO {
	return WithContext(func(ctx context.Context) VIO {
		cb.m.Lock()
		if cb.waiters.Len() == cb.parties.Value()-1 {
			curr := cb.waiters.Front()
			for curr != nil {
				cb.waiters.Remove(curr)
				close(curr.Value)
				curr = cb.waiters.Front()
			}
			cb.m.Unlock()
			return Void()
		}

		waitingChan := make(chan any)
		waiter := cb.waiters.PushBack(waitingChan)
		cb.m.Unlock()

		select {
		case <-ctx.Done():
			cb.m.Lock()
			cb.waiters.Remove(waiter)
			cb.m.Unlock()
			return exitError[v.Void](makeCancellationCause())
		case <-waitingChan:
			return Void()
		}
	})
}

/*
// Below is a more functional approach to a CyclicBarrier. It is not the one used for performance reasons (twice as slower).

func MakeCyclicBarrier(parties r.PosInt) IO[CyclicBarrier] {
	return FlatMap(
		MakeDeferred[v.Void](),
		func(barrier Deferred[v.Void]) IO[CyclicBarrier] {
			return Map(
				MakeRef(barrierState{parties.Value(), 0, barrier}),
				func(ref Ref[barrierState]) CyclicBarrier { return &cyclicBarrier{r: ref, capacity: parties.Value()} },
			)
		},
	)
}

type cyclicBarrier struct {
	r        Ref[barrierState]
	capacity int
}

type barrierState struct {
	awaiting int
	epoch    uint64
	barrier  Deferred[v.Void]
}

func (cb *cyclicBarrier) Await() VIO {
	return FlatMap(MakeDeferred[v.Void](), func(newGate Deferred[v.Void]) IO[v.Void] {
		return Flatten(
			Uncancelable(
				ModifyRef(cb.r, func(s barrierState) (barrierState, VIO) {
					newAwaiting := s.awaiting - 1
					if newAwaiting == 0 {
						return barrierState{cb.capacity, s.epoch + 1, newGate}, s.barrier.Complete(v.Void{}).Void()
					}

					newState := barrierState{newAwaiting, s.epoch, s.barrier}

					restoreCapacity := cb.r.Update(func(b barrierState) barrierState {
						if b.epoch != s.epoch {
							return b
						}
						return barrierState{
							awaiting: b.awaiting - 1,
							barrier:  b.barrier,
						}
					})

					return newState, OnCancelled(s.barrier.Get(), restoreCapacity)
				}),
			),
		)
	})
*/
