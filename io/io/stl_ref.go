package io

import (
	v "github.com/valentinHenry/giog/void"
	"sync"
)

type Ref[A any] interface {
	// Get returns the referenced value
	Get() IO[A]

	// Set sets a new referenced value
	Set(A) VIO

	// Update sets a new value as the updated version of the previous one
	Update(func(A) A) VIO

	// GetAndSet sets a new value and returns the previous one
	GetAndSet(A) IO[A]

	// GetAndUpdate updates the value and returns the previous value
	GetAndUpdate(func(A) A) IO[A]

	// UpdateAndGet updates the value and returns the updated value
	UpdateAndGet(func(A) A) IO[A]
}

func MakeRef[A any](v A) IO[Ref[A]] {
	return Delay(func() Ref[A] {
		var m sync.RWMutex
		return ref[A]{&m, &v}
	})
}

type ref[A any] struct {
	m     *sync.RWMutex
	value *A
}

func (r ref[A]) Get() IO[A] {
	return Delay(func() A {
		r.m.RLock()
		defer r.m.RUnlock()
		return *r.value
	})
}

func (r ref[A]) Set(a A) VIO {
	return Delay(func() v.Void {
		r.m.Lock()
		*r.value = a
		r.m.Unlock()
		return v.Void{}
	})
}

func (r ref[A]) Update(fn func(A) A) VIO {
	return Delay(func() v.Void {
		r.m.Lock()
		*r.value = fn(*r.value)
		r.m.Unlock()
		return v.Void{}
	})
}

func (r ref[A]) GetAndSet(v A) IO[A] {
	return Delay(func() A {
		r.m.Lock()
		value := *r.value
		*r.value = v
		r.m.Unlock()
		return value
	})
}

func (r ref[A]) GetAndUpdate(fn func(A) A) IO[A] {
	return Delay(func() A {
		r.m.Lock()
		value := *r.value
		*r.value = fn(*r.value)
		r.m.Unlock()
		return value
	})
}

func (r ref[A]) UpdateAndGet(fn func(A) A) IO[A] {
	return Delay(func() A {
		r.m.Lock()
		updated := fn(*r.value)
		*r.value = updated
		r.m.Unlock()
		return updated
	})
}
