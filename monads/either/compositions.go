package either

func Compose1[A, B, LeftT any](a Either[LeftT, A], fn1 func(A) Either[LeftT, B]) Either[LeftT, B] {
	return FlatMap(a, fn1)
}

func Compose2[A, B, C, LeftT any](a Either[LeftT, A], fn1 func(A) Either[LeftT, B], fn2 func(B) Either[LeftT, C]) Either[LeftT, C] {
	return FlatMap(Compose1(a, fn1), fn2)
}

func Compose3[A, B, C, D, LeftT any](a Either[LeftT, A], fn1 func(A) Either[LeftT, B], fn2 func(B) Either[LeftT, C], fn3 func(C) Either[LeftT, D]) Either[LeftT, D] {
	return FlatMap(Compose2(a, fn1, fn2), fn3)
}

func Compose4[A, B, C, D, E, LeftT any](a Either[LeftT, A], fn1 func(A) Either[LeftT, B], fn2 func(B) Either[LeftT, C], fn3 func(C) Either[LeftT, D], fn4 func(D) Either[LeftT, E]) Either[LeftT, E] {
	return FlatMap(Compose3(a, fn1, fn2, fn3), fn4)
}

func Compose5[A, B, C, D, E, F, LeftT any](a Either[LeftT, A], fn1 func(A) Either[LeftT, B], fn2 func(B) Either[LeftT, C], fn3 func(C) Either[LeftT, D], fn4 func(D) Either[LeftT, E], fn5 func(E) Either[LeftT, F]) Either[LeftT, F] {
	return FlatMap(Compose4(a, fn1, fn2, fn3, fn4), fn5)
}
