package either

import (
	fn "github.com/valentinHenry/giog/functions"
	o "github.com/valentinHenry/giog/monads/option"
	"github.com/valentinHenry/giog/tuples"
	v "github.com/valentinHenry/giog/void"
)

type Either[L, R any] interface {
	IsLeft() bool
	IsRight() bool
	ToOption() o.Option[R]
	Void() Either[L, v.Void]
	Swap() Either[R, L]

	either()
}

type Left[L, R any] struct {
	LeftValue L
}

func (l Left[L, R]) IsLeft() bool            { return true }
func (l Left[L, R]) IsRight() bool           { return false }
func (l Left[L, R]) ToOption() o.Option[R]   { return o.None[R]{} }
func (l Left[L, R]) Void() Either[L, v.Void] { return Left[L, v.Void]{l.LeftValue} }
func (l Left[L, R]) Swap() Either[R, L]      { return Right[R, L]{l.LeftValue} }

func (l Left[L, R]) either() {}

type Right[L, R any] struct {
	RightValue R
}

func (r Right[L, R]) IsLeft() bool            { return true }
func (r Right[L, R]) IsRight() bool           { return false }
func (r Right[L, R]) ToOption() o.Option[R]   { return o.Some[R]{r.RightValue} }
func (r Right[L, R]) Void() Either[L, v.Void] { return Right[L, v.Void]{v.Void{}} }
func (r Right[L, R]) Swap() Either[R, L]      { return Left[R, L]{r.RightValue} }

func (r Right[L, R]) either() {}

func FromResult[A any](e error, a A) Either[error, A] {
	if e != nil {
		return Left[error, A]{e}
	} else {
		return Right[error, A]{a}
	}
}

func FromOption[L, R any](opt o.Option[R], ifEmpty L) Either[L, R] {
	return o.Fold[R, Either[L, R]](
		opt,
		Left[L, R]{ifEmpty}, func(r R) Either[L, R] { return Right[L, R]{r} },
	)
}

func ToLeft[L, R any](l L) Either[L, R] {
	return Left[L, R]{l}
}

func ToRight[L, R any](r R) Either[L, R] {
	return Right[L, R]{r}
}

func Map[L, R, NR any](e Either[L, R], mapFn func(R) NR) Either[L, NR] {
	return FullMap(e, fn.Identity[L], mapFn)
}
func MapK[L, R, NR any](mapFn func(R) NR) func(Either[L, R]) Either[L, NR] {
	return FullMapK[L, R, L, NR](fn.Identity[L], mapFn)
}

func LeftMap[L, R, NL any](e Either[L, R], mapFn func(L) NL) Either[NL, R] {
	return FullMap(e, mapFn, fn.Identity[R])
}
func LeftMapK[L, R, NL any](mapFn func(L) NL) func(Either[L, R]) Either[NL, R] {
	return FullMapK[L, R, NL, R](mapFn, fn.Identity[R])
}

func ToOption[L, R any](e Either[L, R]) o.Option[R] {
	return Fold(
		e,
		func(l L) o.Option[R] { return o.None[R]{} },
		func(r R) o.Option[R] { return o.Some[R]{r} },
	)
}
func ToOptionK[L, R any]() func(Either[L, R]) o.Option[R] {
	return func(e Either[L, R]) o.Option[R] {
		return Fold(
			e,
			func(l L) o.Option[R] { return o.None[R]{} },
			func(r R) o.Option[R] { return o.Some[R]{r} },
		)
	}
}

func FlatMap[L, R, NR any](e Either[L, R], mapFn func(R) Either[L, NR]) Either[L, NR] {
	switch e := e.(type) {
	case Left[L, R]:
		return Left[L, NR]{e.LeftValue}
	case Right[L, R]:
		return mapFn(e.RightValue)
	default:
		return nil
	}
}
func FlatMapK[L, R, NR any](mapFn func(R) Either[L, NR]) func(Either[L, R]) Either[L, NR] {
	return func(e Either[L, R]) Either[L, NR] { return FlatMap(e, mapFn) }
}

func FullMap[L, R, NL, NR any](e Either[L, R], onLeft func(L) NL, onRight func(R) NR) Either[NL, NR] {
	switch e := e.(type) {
	case Left[L, R]:
		return Left[NL, NR]{onLeft(e.LeftValue)}
	case Right[L, R]:
		return Right[NL, NR]{onRight(e.RightValue)}
	default:
		return nil
	}
}
func FullMapK[L, R, NL, NR any](onLeft func(L) NL, onRight func(R) NR) func(Either[L, R]) Either[NL, NR] {
	return func(e Either[L, R]) Either[NL, NR] { return FullMap(e, onLeft, onRight) }
}

func Fold[L, R, O any](e Either[L, R], onLeft func(L) O, onRight func(R) O) O {
	switch e := e.(type) {
	case Left[L, R]:
		return onLeft(e.LeftValue)
	case Right[L, R]:
		return onRight(e.RightValue)
	default:
		var nilO O
		return nilO
	}
}
func FoldK[L, R, O any](onLeft func(L) O, onRight func(R) O) func(Either[L, R]) O {
	return func(e Either[L, R]) O { return Fold(e, onLeft, onRight) }
}

func JoinRight[E, R1, R2 any](l Either[E, R1], r Either[E, R2]) Either[E, tuples.T2[R1, R2]] {
	return LeftMap(Join(l, r), Merge[E])
}
func JoinRightK[E, R1, R2 any](r Either[E, R2]) func(l Either[E, R1]) Either[E, tuples.T2[R1, R2]] {
	return func(l Either[E, R1]) Either[E, tuples.T2[R1, R2]] { return JoinRight(l, r) }
}

func Join[E1, E2, R1, R2 any](l Either[E1, R1], r Either[E2, R2]) Either[Either[E1, E2], tuples.T2[R1, R2]] {
	mr := LeftMap(r, func(e2 E2) Either[E1, E2] { return Right[E1, E2]{e2} })

	return Fold(
		l,
		func(e E1) Either[Either[E1, E2], tuples.T2[R1, R2]] {
			return Left[Either[E1, E2], tuples.T2[R1, R2]]{Left[E1, E2]{e}}
		},
		func(r1 R1) Either[Either[E1, E2], tuples.T2[R1, R2]] {
			return Map(
				mr,
				func(r2 R2) tuples.T2[R1, R2] {
					return tuples.Of2(r1, r2)
				},
			)
		},
	)
}
func JoinK[E1, E2, R1, R2 any](r Either[E2, R2]) func(Either[E1, R1]) Either[Either[E1, E2], tuples.T2[R1, R2]] {
	return func(l Either[E1, R1]) Either[Either[E1, E2], tuples.T2[R1, R2]] {
		return Join(l, r)
	}
}

func Merge[A any](e Either[A, A]) A {
	return Fold(e, fn.Identity[A], fn.Identity[A])
}
func MergeK[A any]() func(Either[A, A]) A {
	return func(e Either[A, A]) A { return Merge(e) }
}
