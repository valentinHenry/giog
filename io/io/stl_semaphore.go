package io

type Semaphore interface {
	Available() IO[uint]
	Acquire() VIO
	Release() VIO
}
