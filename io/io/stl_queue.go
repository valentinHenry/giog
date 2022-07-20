package io

import (
	"context"
	dll "github.com/valentinHenry/giog/utils/collections/double_link_list"
	o "github.com/valentinHenry/giog/utils/monads/option"
	v "github.com/valentinHenry/giog/utils/void"
	r "github.com/valentinHenry/refined"
	"math"
	"sync"
)

type Queue[A any] interface {
	Enqueue(a A) IO[v.Void]
	TryEnqueue(a A) IO[bool]

	Dequeue() IO[A]
	TryDequeue() IO[o.Option[A]]

	Head() IO[A]
	TryHead() IO[o.Option[A]]
}

func BoundedQueue[A any](nb r.PosInt) IO[Queue[A]] {
	return Delay(func() Queue[A] {
		var m sync.Mutex
		elts := dll.New[A]()
		waitersEnqueue := dll.New[enqueueWaiter[A]]()
		waitersDequeue := dll.New[dequeueWaiter[A]]()
		return queue[A]{
			m:              &m,
			maxLength:      nb.Value(),
			elts:           elts,
			waitersEnqueue: waitersEnqueue,
			waitersDequeue: waitersDequeue,
		}
	})
}

func UnboundedQueue[A any]() IO[Queue[A]] {
	return Delay(func() Queue[A] {
		var m sync.Mutex
		elts := dll.New[A]()
		waitersEnqueue := dll.New[enqueueWaiter[A]]()
		waitersDequeue := dll.New[dequeueWaiter[A]]()
		return queue[A]{
			m:              &m,
			maxLength:      math.MaxInt,
			elts:           elts,
			waitersEnqueue: waitersEnqueue,
			waitersDequeue: waitersDequeue,
		}
	})
}

type queue[A any] struct {
	m              *sync.Mutex
	maxLength      int
	elts           *dll.List[A]
	waitersEnqueue *dll.List[enqueueWaiter[A]]
	waitersDequeue *dll.List[dequeueWaiter[A]]
}

func (q queue[A]) Enqueue(a A) VIO {
	// FIXME use Blocking instead of blocking select with ctx usage
	return WithContext(func(ctx context.Context) VIO {
		return Lift(func() (v.Void, error) {
			q.m.Lock()
			if q.elts.Len() < q.maxLength {
				q.elts.PushBack(a)
				q.notifyDequeueWaiters()
				q.m.Unlock()
				return v.Void{}, nil
			}

			ready := make(chan struct{})
			w := enqueueWaiter[A]{
				ready: ready,
				value: a,
			}
			elem := q.waitersEnqueue.PushBack(w)
			q.m.Unlock()

			select {
			case <-ctx.Done():
				q.m.Lock()
				q.waitersEnqueue.Remove(elem)
				q.m.Unlock()
				return v.Void{}, ctx.Err()
			case <-ready:
				return v.Void{}, nil
			}
		})
	})
}

func (q queue[A]) TryEnqueue(a A) IO[bool] {
	return Delay(func() bool {
		q.m.Lock()

		res := q.elts.Len() < q.maxLength
		if res {
			q.elts.PushBack(a)
			q.notifyDequeueWaiters()
		}
		q.m.Unlock()
		return res
	})
}

func (q queue[A]) Dequeue() IO[A] {
	// FIXME use Blocking instead of blocking select with ctx usage
	return WithContext(func(ctx context.Context) IO[A] {
		return Lift(func() (A, error) {
			q.m.Lock()
			if q.elts.Len() > 0 {
				f := q.elts.Front()
				q.elts.Remove(f)
				q.notifyEnqueueWaiters()
				q.m.Unlock()
				return f.Value, nil
			}

			ready := make(chan A)
			w := dequeueWaiter[A]{
				ready: ready,
			}
			elem := q.waitersDequeue.PushBack(w)
			q.m.Unlock()

			select {
			case <-ctx.Done():
				q.m.Lock()
				q.waitersDequeue.Remove(elem)
				q.m.Unlock()
				var dummyA A
				return dummyA, ctx.Err()
			case dequeued := <-ready:
				return dequeued, nil
			}
		})
	})
}

func (q queue[A]) TryDequeue() IO[o.Option[A]] {
	return Delay(func() o.Option[A] {
		q.m.Lock()
		defer q.m.Unlock()

		if q.elts.Len() == 0 {
			return o.None[A]{}
		}
		head := q.elts.Front()
		q.elts.Remove(head)
		q.notifyEnqueueWaiters()
		return o.Some[A]{head.Value}
	})
}

func (q queue[A]) Head() IO[A] {
	// TODO
	return _TODO[IO[A]]()
}

func (q queue[A]) TryHead() IO[o.Option[A]] {
	// TODO
	return _TODO[IO[o.Option[A]]]()
}

func (q queue[A]) notifyEnqueueWaiters() {
	for {
		next := q.waitersEnqueue.Front()

		if next == nil {
			break
		}

		if q.elts.Len() >= q.maxLength {
			break
		}

		value := next.Value.value
		channel := next.Value.ready

		q.elts.PushBack(value)
		q.waitersEnqueue.Remove(next)
		close(channel)
	}
}

func (q queue[A]) notifyDequeueWaiters() {
	for {
		next := q.waitersDequeue.Front()

		if next == nil {
			break
		}

		if q.elts.Len() == 0 {
			break
		}

		channel := next.Value.ready

		head := q.elts.Front()
		q.elts.Remove(head)
		q.waitersDequeue.Remove(next)
		channel <- head.Value
		close(channel)
	}
}

type enqueueWaiter[A any] struct {
	ready chan<- struct{}
	value A
}

type dequeueWaiter[A any] struct {
	ready chan<- A
}
