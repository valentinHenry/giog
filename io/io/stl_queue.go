package io

import (
	"context"
	dll "github.com/valentinHenry/giog/utils/collections/double_link_list"
	o "github.com/valentinHenry/giog/utils/monads/option"
	v "github.com/valentinHenry/giog/utils/void"
	"math"
	"sync"
)

// Queue is an interface to a concurrent queue
//
// Three implementations are available:
//
// - BoundedQueue:
// A queue which, when full, will block semantically on Enqueue(A) and return
// false on TryEnqueue(A)
//
// - UnboundedQueue:
// A limitless queue
//
// - SyncQueue
// A queue which requires at least one reader to enqueue a value.
type Queue[A any] interface {
	// Enqueue enqueues a value, semantically blocks in case there is no available
	// place.
	Enqueue(a A) IO[v.Void]

	// TryEnqueue enqueues a value if the queue is not full and returns true, it
	// returns false otherwise.
	TryEnqueue(a A) IO[bool]

	// Dequeue dequeues a value if there is one in the queue, it blocks otherwise
	// until one is available.
	Dequeue() IO[A]

	// TryDequeue dequeues a value if there is one in the queue and returns a
	// defined Option with the dequeued value, it returns None otherwise.
	TryDequeue() IO[o.Option[A]]
}

func BoundedQueue[A any](nb uint) IO[Queue[A]] {
	return Delay(func() Queue[A] {
		elts := dll.New[A]()
		waitersEnqueue := dll.New[enqueueWaiter[A]]()
		waitersDequeue := dll.New[dequeueWaiter[A]]()
		return &queue[A]{
			m:              sync.Mutex{},
			maxLength:      nb,
			elts:           elts,
			enqueueWaiters: waitersEnqueue,
			dequeueWaiters: waitersDequeue,
		}
	})
}

func UnboundedQueue[A any]() IO[Queue[A]] {
	return Delay(func() Queue[A] {
		elts := dll.New[A]()
		waitersEnqueue := dll.New[enqueueWaiter[A]]()
		waitersDequeue := dll.New[dequeueWaiter[A]]()
		return &queue[A]{
			m:              sync.Mutex{},
			maxLength:      math.MaxUint,
			elts:           elts,
			enqueueWaiters: waitersEnqueue,
			dequeueWaiters: waitersDequeue,
		}
	})
}

func SyncQueue[A any]() IO[Queue[A]] {
	return Delay(func() Queue[A] {
		elts := dll.New[A]()
		waitersEnqueue := dll.New[enqueueWaiter[A]]()
		waitersDequeue := dll.New[dequeueWaiter[A]]()
		return &queue[A]{
			m:              sync.Mutex{},
			maxLength:      0,
			elts:           elts,
			enqueueWaiters: waitersEnqueue,
			dequeueWaiters: waitersDequeue,
		}
	})
}

type queue[A any] struct {
	m              sync.Mutex
	maxLength      uint
	elts           *dll.List[A]
	enqueueWaiters *dll.List[enqueueWaiter[A]]
	dequeueWaiters *dll.List[dequeueWaiter[A]]
}

func (q *queue[A]) Enqueue(a A) VIO {
	return WithContext(func(ctx context.Context) VIO {
		q.m.Lock()
		if q.elts.Len() < q.maxLength+q.dequeueWaiters.Len() {
			q.elts.PushBack(a)
			q.notifyNextDequeueWaiter()
			q.m.Unlock()
			return Void()
		}

		ready := make(chan struct{})
		w := enqueueWaiter[A]{
			ready: ready,
			value: a,
		}
		elem := q.enqueueWaiters.PushBack(w)
		q.m.Unlock()

		select {
		case <-ctx.Done():
			q.m.Lock()
			select {
			case <-ready:
				q.notifyNextDequeueWaiter()
			default:
				q.enqueueWaiters.Remove(elem)
			}
			q.m.Unlock()
			return exitError[v.Void](makeCancellationCause())
		case <-ready:
			q.m.Lock()
			q.notifyNextDequeueWaiter()
			q.m.Unlock()
			return Void()
		}
	})
}

func (q *queue[A]) TryEnqueue(a A) IO[bool] {
	return Delay(func() bool {
		q.m.Lock()

		res := q.elts.Len() < q.maxLength+q.dequeueWaiters.Len()
		if res {
			q.elts.PushBack(a)
			q.notifyNextDequeueWaiter()
		}
		q.m.Unlock()
		return res
	})
}

func (q *queue[A]) Dequeue() IO[A] {
	return WithContext(func(ctx context.Context) IO[A] {
		q.m.Lock()
		if int64(q.elts.Len()) > -int64(q.enqueueWaiters.Len()) {
			q.notifyNextEnqueueWaiter()
			f := q.elts.Front()
			q.elts.Remove(f)
			q.m.Unlock()
			return Pure(f.Value)
		}

		ready := make(chan A)
		w := dequeueWaiter[A]{
			ready: ready,
		}
		elem := q.dequeueWaiters.PushBack(w)
		q.m.Unlock()

		select {
		case <-ctx.Done():
			q.m.Lock()
			select {
			case <-ready:
				q.notifyNextEnqueueWaiter()
			default:
				q.dequeueWaiters.Remove(elem)
			}
			q.m.Unlock()
			return exitError[A](makeCancellationCause())
		case dequeued := <-ready:
			q.m.Lock()
			q.notifyNextEnqueueWaiter()
			q.m.Unlock()
			return Pure(dequeued)
		}
	})
}

func (q *queue[A]) TryDequeue() IO[o.Option[A]] {
	return Delay(func() o.Option[A] {
		q.m.Lock()
		defer q.m.Unlock()

		if q.elts.Len()+q.enqueueWaiters.Len() == 0 {
			return o.None[A]{}
		}

		q.notifyNextEnqueueWaiter()
		head := q.elts.Front()
		q.elts.Remove(head)
		return o.Some[A]{head.Value}
	})
}

type enqueueWaiter[A any] struct {
	ready chan<- struct{}
	value A
}

type dequeueWaiter[A any] struct {
	ready chan<- A
}

func (q *queue[A]) notifyNextDequeueWaiter() {
	waiter := q.dequeueWaiters.Front()
	if waiter == nil {
		return
	}
	value := q.elts.Front()
	waiter.Value.ready <- value.Value

	q.dequeueWaiters.Remove(waiter)
	q.elts.Remove(value)
}

func (q *queue[A]) notifyNextEnqueueWaiter() {
	waiter := q.enqueueWaiters.Front()
	if waiter == nil {
		return
	}
	q.elts.PushBack(waiter.Value.value)
	close(waiter.Value.ready)
}
