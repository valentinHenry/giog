package io

import (
	"context"
	dll "github.com/valentinHenry/giog/utils/collections/double_link_list"
	v "github.com/valentinHenry/giog/utils/void"
	"sync"
)

// Semaphore is an interface to a non-negative amount of permits.
type Semaphore interface {
	Release() VIO
	ReleaseN(n uint) VIO

	Acquire() VIO
	AcquireN(n uint) VIO

	TryAcquire() IO[bool]
	TryAcquireN(n uint) IO[bool]

	Use() RIO[v.Void]
	UseN(n uint) RIO[v.Void]
}

// MakeSemaphore is an effect returning an instance of a Semaphore.
// size is the number of permits available in the semaphore
func MakeSemaphore(size uint) IO[Semaphore] {
	return Delay[Semaphore](func() Semaphore {
		return &semaphore{
			waiters:          dll.New[waiter](),
			m:                sync.Mutex{},
			availablePermits: size,
			size:             size,
		}
	})
}

type semaphore struct {
	waiters          *dll.List[waiter]
	m                sync.Mutex
	availablePermits uint
	size             uint
}

func (s *semaphore) Release() VIO {
	return s.ReleaseN(1)
}
func (s *semaphore) ReleaseN(n uint) VIO {
	return Defer(func() VIO {
		s.m.Lock()
		defer s.m.Unlock()

		if s.availablePermits+n > s.size {
			return Raise[v.Void](IllegalReleasing{})
		}

		s.availablePermits -= n
		s.notifyWaiters()
		return Void()
	})
}

func (s *semaphore) Acquire() VIO {
	return s.AcquireN(1)
}

func (s *semaphore) AcquireN(n uint) VIO {
	return WithContext(func(ctx context.Context) VIO {
		s.m.Lock()
		if s.availablePermits >= n {
			s.availablePermits -= n
			s.m.Unlock()
			return Void()
		}

		wchan := make(chan any)
		w := s.waiters.PushBack(waiter{wchan, n})

		s.m.Unlock()
		select {
		case <-ctx.Done():
			s.m.Lock()
			select {
			case <-wchan:
				s.availablePermits += n
				s.notifyWaiters()
			default:
				s.waiters.Remove(w)
			}
			s.m.Unlock()
			return exitError[v.Void](makeCancellationCause())
		case <-wchan:
			return Void()
		}
	})
}

func (s *semaphore) TryAcquire() IO[bool] {
	return s.TryAcquireN(1)
}

func (s *semaphore) TryAcquireN(n uint) IO[bool] {
	return Delay(func() bool {
		s.m.Lock()
		defer s.m.Unlock()

		if s.availablePermitsAfterWaitersUnsafe() >= n {
			s.availablePermits -= n
			return true
		}

		return false
	})
}

func (s *semaphore) Use() RIO[v.Void] {
	return s.UseN(1)
}

func (s *semaphore) UseN(n uint) RIO[v.Void] {
	return MakeRIO(s.AcquireN(n), func(_ v.Void) VIO { return s.ReleaseN(n) })
}

type IllegalReleasing struct{}

func (i IllegalReleasing) Error() string {
	return "illegal releasing of semaphore: not enough permits available"
}

type waiter struct {
	wait    chan any
	permits uint
}

func (s *semaphore) availablePermitsAfterWaitersUnsafe() uint {
	available := s.availablePermits
	if s.waiters.Len() > 0 {
		curr := s.waiters.Front().Next()
		for ; curr != s.waiters.Front(); curr = curr.Next() {
			if available < curr.Value.permits {
				return 0
			}
			available -= curr.Value.permits
		}
	}
	return available
}

func (s *semaphore) notifyWaiters() {
	for curr := s.waiters.Front(); curr != nil; curr = s.waiters.Front() {
		if s.availablePermits < curr.Value.permits {
			return
		}
		s.availablePermits -= curr.Value.permits
		close(curr.Value.wait)
		s.waiters.Remove(curr)
	}
}
