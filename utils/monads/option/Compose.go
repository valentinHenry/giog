package option

func Compose1[A, B any](a Option[A], fn1 func(A) Option[B]) Option[B] {
	return FlatMap(a, fn1)
}

func Compose2[A, B, C any](a Option[A], fn1 func(A) Option[B], fn2 func(B) Option[C]) Option[C] {
	return FlatMap(Compose1(a, fn1), fn2)
}

func Compose3[A, B, C, D any](a Option[A], fn1 func(A) Option[B], fn2 func(B) Option[C], fn3 func(C) Option[D]) Option[D] {
	return FlatMap(Compose2(a, fn1, fn2), fn3)
}

func Compose4[A, B, C, D, E any](a Option[A], fn1 func(A) Option[B], fn2 func(B) Option[C], fn3 func(C) Option[D], fn4 func(D) Option[E]) Option[E] {
	return FlatMap(Compose3(a, fn1, fn2, fn3), fn4)
}

func Compose5[A, B, C, D, E, F any](a Option[A], fn1 func(A) Option[B], fn2 func(B) Option[C], fn3 func(C) Option[D], fn4 func(D) Option[E], fn5 func(E) Option[F]) Option[F] {
	return FlatMap(Compose4(a, fn1, fn2, fn3, fn4), fn5)
}

func CMap[A, B any](fn func(A) B) func(A) Option[B] {
	return func(a A) Option[B] {
		return Some[B]{fn(a)}
	}
}

func CTap[A any](fn func(A)) func(A) Option[A] {
	return func(a A) Option[A] {
		fn(a)
		return Some[A]{a}
	}
}
