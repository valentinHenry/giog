package functions

import (
	o "github.com/valentinHenry/giog/monads/option"
)

type Product[A, B any] interface{}

type PartialFunction[A, B any] interface {
	Apply(any) o.Option[B]
}

type _PartialFn[A, B any] func(A) B

func (p _PartialFn[A, B]) Apply(a any) o.Option[B] {
	tpe, ok := a.(A)
	return o.When(ok, p(tpe))
}

type _ConditionalPartialFn[A, B any] struct {
	condition func(A) bool
	fn        func(A) B
}

func (c _ConditionalPartialFn[A, B]) Apply(a any) o.Option[B] {
	tpe, ok := a.(A)
	return o.When(ok && c.condition(tpe), c.fn(tpe))
}

type _CombinedPartialFn[A1, A2, B any] struct {
	first  PartialFunction[A1, B]
	second PartialFunction[A2, B]
}

func (c _CombinedPartialFn[A1, A2, B]) Apply(a any) o.Option[B] {
	res := c.first.Apply(a)
	if res.IsDefined() {
		return res
	}
	return c.second.Apply(a)
}

func ParFunc[A, B any](fn func(A) B) PartialFunction[A, B] {
	var parFn _PartialFn[A, B] = fn
	return parFn
}

func ParCondFunc[A, B any](cond func(A) bool, fn func(A) B) PartialFunction[A, B] {
	return _ConditionalPartialFn[A, B]{
		condition: cond,
		fn:        fn,
	}
}

func Combine[A1, A2, NA, B any](first PartialFunction[A1, B], second PartialFunction[A2, B]) PartialFunction[NA, B] {
	return _CombinedPartialFn[A1, A2, B]{
		first:  first,
		second: second,
	}
}
