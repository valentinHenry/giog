package io

import o "github.com/valentinHenry/giog/utils/monads/option"

type OptionT[A any] IO[o.Option[A]]

// Semi (IO)
// Sub (Option)

func OTSubFold[A, B any](io IO[o.Option[A]], ifEmpty B, ifDefined func(A) B) IO[B] {
	return _OTSubFold(getTrace(1), io, ifEmpty, ifDefined)
}
func _OTSubFold[A, B any](_trace *Trace, io IO[o.Option[A]], ifEmpty B, ifDefined func(A) B) IO[B] {
	return _Map(_trace, io, func(v o.Option[A]) B { return o.Fold(v, ifEmpty, ifDefined) })
}

func OTSubFoldIO[A, B any](io IO[o.Option[A]], ifEmpty IO[B], ifDefined func(A) IO[B]) IO[B] {
	return _OTSubFoldIO(getTrace(1), io, ifEmpty, ifDefined)
}
func _OTSubFoldIO[A, B any](_trace *Trace, io IO[o.Option[A]], ifEmpty IO[B], ifDefined func(A) IO[B]) IO[B] {
	return _FlatMap(_trace, io, func(v o.Option[A]) IO[B] { return o.Fold(v, ifEmpty, ifDefined) })
}
