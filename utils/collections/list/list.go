package list

import o "github.com/valentinHenry/giog/utils/monads/option"

type List[T any] interface {
	Append(v T) List[T]
	HeadOpt() o.Option[T]
	TailOpt() o.Option[List[T]]
}

type Cons[T any] struct {
	Head T
	Tail List[T]
}

type Nil[T any] struct{}

func (c *Cons[T]) Append(v T) List[T] {
	return &Cons[T]{Head: v, Tail: c}
}
func (c *Cons[T]) HeadOpt() o.Option[T] {
	return o.Some[T]{c.Head}
}
func (c *Cons[T]) TailOpt() o.Option[List[T]] {
	return o.Some[List[T]]{c.Tail}
}

func (n *Nil[T]) Append(v T) List[T] {
	return &Cons[T]{Head: v, Tail: n}
}
func (n *Nil[T]) HeadOpt() o.Option[T] {
	return o.None[T]{}
}
func (n *Nil[T]) TailOpt() o.Option[List[T]] {
	return o.None[List[T]]{}
}
