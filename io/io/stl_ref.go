package io

import (
	"errors"
	v "github.com/valentinHenry/giog/utils/void"
	"sync/atomic"
)

// Ref is an interface representing the reference to a value.
type Ref[A any] interface {
	// Get returns the referenced value.
	Get() IO[A]

	// Set sets a new referenced value.
	Set(A) VIO

	// Update sets a new value as the updated version of the previous one.
	Update(func(A) A) VIO

	// TryUpdate tries to set a new value as the updated version of the previous
	// one. It succeeds if no-one is updating the value
	TryUpdate(func(A) A) IO[bool]

	// GetAndSet sets a new value and returns the previous one.
	GetAndSet(A) IO[A]

	// GetAndUpdate updates the value and returns the previous value.
	GetAndUpdate(func(A) A) IO[A]

	// UpdateAndGet updates the value and returns the updated value.
	UpdateAndGet(func(A) A) IO[A]
}

// MakeRef creates a Ref[A] which can be used later on
func MakeRef[A any](v A) IO[Ref[A]] {
	return Delay(func() Ref[A] {
		r := &ref[A]{atomic.Value{}}
		r.value.Store(&v)
		return r
	})
}

type ref[A any] struct {
	value atomic.Value
}

func (r *ref[A]) Get() IO[A] {
	return Delay(func() A {
		return *r.value.Load().(*A)
	})
}

func (r *ref[A]) Set(a A) VIO {
	return Delay(func() v.Void {
		r.value.Store(&a)
		return v.Void{}
	})
}

func (r *ref[A]) Update(fn func(A) A) VIO {
	update := func() v.Void {
		cond := false
		for !cond {
			oldValue := r.value.Load().(*A)
			newValue := fn(*oldValue)
			cond = r.value.CompareAndSwap(oldValue, &newValue)
		}
		return v.Void{}
	}

	return Delay(update)
}

func (r *ref[A]) TryUpdate(fn func(A) A) IO[bool] {
	return Delay(func() bool {
		oldValue := r.value.Load().(*A)
		newValue := fn(*oldValue)
		return r.value.CompareAndSwap(oldValue, &newValue)
	})
}

func (r *ref[A]) GetAndSet(v A) IO[A] {
	return Delay(func() A { return *r.value.Swap(v).(*A) })
}

func (r *ref[A]) GetAndUpdate(fn func(A) A) IO[A] {
	setValue := func() A {
		cond := false
		var oldValue *A

		for !cond {
			oldValue = r.value.Load().(*A)
			newValue := fn(*oldValue)
			cond = r.value.CompareAndSwap(oldValue, &newValue)
		}

		return *oldValue
	}

	return Delay(setValue)
}

func (r *ref[A]) UpdateAndGet(fn func(A) A) IO[A] {
	setValue := func() A {
		cond := false
		var newValue A

		for !cond {
			oldValue := r.value.Load().(*A)
			newValue = fn(*oldValue)
			cond = r.value.CompareAndSwap(oldValue, &newValue)
		}

		return newValue
	}

	return Delay(setValue)
}

func ModifyRef[A, B any](r Ref[A], modify func(A) (A, B)) IO[B] {
	if rf, ok := r.(*ref[A]); ok {
		setValue := func() B {
			cond := false
			var newValue A
			var ret B

			for !cond {
				oldValue := rf.value.Load().(*A)
				newValue, ret = modify(*oldValue)
				cond = rf.value.CompareAndSwap(oldValue, &newValue)
			}

			return ret
		}

		return Delay(setValue)
	} else {
		return _Raise[B](getTrace(1), errors.New("cannot modify another implementation of Ref"))
	}
}
