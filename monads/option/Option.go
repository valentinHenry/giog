package option

import (
	"errors"
	"fmt"
)

type Option[V any] interface {
	IsEmpty() bool
	IsDefined() bool

	GetOrElse(orElse V) V

	value() (V, error)
}

type None[V any] struct{}

func (n None[V]) IsEmpty() bool        { return true }
func (n None[V]) IsDefined() bool      { return false }
func (n None[V]) GetOrElse(orElse V) V { return orElse }
func (n None[V]) value() (V, error) {
	var dummy V
	return dummy, errors.New("tried to get a value from an empty option")
}
func (n None[V]) String() string { return fmt.Sprint("None") }

type Some[V any] struct {
	Value V
}

func (s Some[V]) IsEmpty() bool        { return false }
func (s Some[V]) IsDefined() bool      { return true }
func (s Some[V]) GetOrElse(orElse V) V { return s.Value }
func (s Some[V]) value() (V, error)    { return s.Value, nil }
func (s Some[V]) String() string       { return fmt.Sprint("Some(", s.Value, ")") }

func Map[A, B any](o Option[A], mapFn func(A) B) Option[B] {
	switch o := o.(type) {
	case Some[A]:
		return Some[B]{mapFn(o.Value)}
	case None[A]:
		return None[B]{}
	default:
		panic("Impossible")
	}
}
func MapK[A, B any](mapFn func(A) B) func(Option[A]) Option[B] {
	return func(o Option[A]) Option[B] { return Map(o, mapFn) }
}

func FlatMap[A, B any](o Option[A], mapFn func(A) Option[B]) Option[B] {
	switch o := o.(type) {
	case Some[A]:
		return mapFn(o.Value)
	case None[A]:
		return None[B]{}
	default:
		panic("Impossible")
	}
}

func FlatMapK[A, B any](mapFn func(A) Option[B]) func(Option[A]) Option[B] {
	return func(e Option[A]) Option[B] { return FlatMap(e, mapFn) }
}

func Fold[A, B any](o Option[A], ifEmpty B, ifDefined func(A) B) B {
	switch o := o.(type) {
	case None[A]:
		return ifEmpty
	case Some[A]:
		return ifDefined(o.Value)
	default:
		panic("Impossible")
	}
}
func FoldK[A, B any](ifEmpty B, ifDefined func(A) B) func(Option[A]) B {
	return func(e Option[A]) B { return Fold(e, ifEmpty, ifDefined) }
}

func When[A any](cond bool, ifTrue A) Option[A] {
	if !cond {
		return None[A]{}
	}
	return Some[A]{ifTrue}
}

func GetOrElse[A any](o Option[A], ifEmpty A) A {
	return o.GetOrElse(ifEmpty)
}
func GetOrElseK[A any](ifEmpty A) func(Option[A]) A {
	return func(o Option[A]) A { return o.GetOrElse(ifEmpty) }
}
