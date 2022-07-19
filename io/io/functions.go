package io

import (
	"context"
	f "github.com/valentinHenry/giog/functions"
	e "github.com/valentinHenry/giog/monads/either"
	o "github.com/valentinHenry/giog/monads/option"
	t "github.com/valentinHenry/giog/tuples"
	v "github.com/valentinHenry/giog/void"
	r "github.com/valentinHenry/refined"
	"time"
)

func Map[T1, T2 any](io IO[T1], mapFn func(T1) T2) IO[T2] {
	return _Map(getTrace(1), io, mapFn)
}
func MapK[T1, T2 any](mapFn func(T1) T2) func(IO[T1]) IO[T2] {
	return _MapK(getTrace(1), mapFn)
}

func MapBoth[T1, T2 any](io IO[T1], onError func(error) error, onSuccess func(T1) T2) IO[T2] {
	return _MapBoth(getTrace(1), io, onError, onSuccess)
}
func MapBothK[T1, T2 any](onError func(error) error, onSuccess func(T1) T2) func(IO[T1]) IO[T2] {
	return _MapBothK(getTrace(1), onError, onSuccess)
}

func FlatMap[T1, T2 any](io IO[T1], mapFn func(T1) IO[T2]) IO[T2] {
	return _FlatMap(getTrace(1), io, mapFn)
}
func FlatMapK[T1, T2 any](mapFn func(T1) IO[T2]) func(IO[T1]) IO[T2] {
	return _FlatMapK(getTrace(1), mapFn)
}

func FlatTap[T1, T2 any](io IO[T1], mapFn func(T1) IO[T2]) IO[T1] {
	return _FlatTap(getTrace(1), io, mapFn)
}
func FlatTapK[T1, T2 any](mapFn func(T1) IO[T2]) func(IO[T1]) IO[T1] {
	return _FlatTapK(getTrace(1), mapFn)
}

func Flatten[T any](io IO[IO[T]]) IO[T] {
	return _FlatMap(getTrace(1), io, f.Identity[IO[T]])
}
func FlattenK[T any]() func(io IO[IO[T]]) IO[T] {
	trace := getTrace(1)
	return func(io IO[IO[T]]) IO[T] {
		return _FlatMap(trace, io, f.Identity[IO[T]])
	}
}

func As[T1, T2 any](io IO[T1], v T2) IO[T2] {
	return _As(getTrace(1), io, v)
}
func AsK[T1, T2 any](v T2) func(IO[T1]) IO[T2] {
	return _AsK[T1, T2](getTrace(1), v)
}

func Absolve[T any](io func() IO[e.Either[error, T]]) IO[T] {
	return _Absolve(getTrace(1), io)
}
func AbsolveK[T any]() func(func() IO[e.Either[error, T]]) IO[T] {
	return _AbsolveK[T](getTrace(1))
}

func FromEither[T any](either func() e.Either[error, T]) IO[T] {
	return _FromEither(getTrace(1), either)
}
func FromEitherK[T any]() func(func() e.Either[error, T]) IO[T] {
	return _FromEitherK[T](getTrace(1))
}

func Pure[T any](v T) IO[T] {
	return _Pure(v)
}
func PureK[T any]() func(T) IO[T] {
	return _PureK[T]()
}

func Raise[T any](err error) IO[T] {
	return _Raise[T](getTrace(1), err)
}
func RaiseK[T any]() func(error) IO[T] {
	return _RaiseK[T](getTrace(1))
}

func Delay[T any](v func() T) IO[T] {
	return _Delay(getTrace(1), v)
}
func DelayK[T any]() func(func() T) IO[T] {
	return _DelayK[T](getTrace(1))
}

func Defer[T any](io func() IO[T]) IO[T] {
	return _Defer(getTrace(1), io)
}
func DeferK[T any]() func(func() IO[T]) IO[T] {
	return _DeferK[T](getTrace(1))
}

func Lift[T any](v func() (T, error)) IO[T] {
	return _Lift(getTrace(1), v)
}
func LiftK[T any]() func(func() (T, error)) IO[T] {
	return _LiftK[T](getTrace(1))
}

func LiftV(v func() error) VIO {
	return _LiftV(getTrace(1), v)
}
func LiftVK() func(func() error) VIO {
	return _LiftVK(getTrace(1))
}

func Fold[T1, T2 any](io IO[T1], onError func(error) T2, onSuccess func(T1) T2) IO[T2] {
	return _Fold(getTrace(1), io, onError, onSuccess)
}
func FoldK[T1, T2 any](onError func(error) T2, onSuccess func(T1) T2) func(IO[T1]) IO[T2] {
	return _FoldK[T1, T2](getTrace(1), onError, onSuccess)
}

func FoldIO[T1, T2 any](io IO[T1], onError func(error) IO[T2], onSuccess func(T1) IO[T2]) IO[T2] {
	return _FoldIO(getTrace(1), io, onError, onSuccess)
}
func FoldIOK[T1, T2 any](onError func(error) IO[T2], onSuccess func(T1) IO[T2]) func(IO[T1]) IO[T2] {
	return _FoldIOK[T1, T2](getTrace(1), onError, onSuccess)
}

func FoldCause[T1, T2 any](io IO[T1], onError func(Cause) T2, onSuccess func(T1) T2) IO[T2] {
	return _FoldCause(getTrace(1), io, onError, onSuccess)
}
func FoldCauseK[T1, T2 any](onError func(Cause) T2, onSuccess func(T1) T2) func(IO[T1]) IO[T2] {
	return _FoldCauseK(getTrace(1), onError, onSuccess)
}

func FoldIOCause[T1, T2 any](io IO[T1], onError func(Cause) IO[T2], onSuccess func(T1) IO[T2]) IO[T2] {
	return _FoldIOCause(getTrace(1), io, onError, onSuccess)
}
func FoldIOCauseK[T1, T2 any](onError func(Cause) IO[T2], onSuccess func(T1) IO[T2]) func(IO[T1]) IO[T2] {
	return _FoldIOCauseK(getTrace(1), onError, onSuccess)
}

func OnSuccess[T1, T2 any](io IO[T1], fn func(T1) T2) IO[o.Option[T2]] {
	return _OnSuccess(getTrace(1), io, fn)
}
func OnSuccessK[T1, T2 any](fn func(T1) T2) func(IO[T1]) IO[o.Option[T2]] {
	return _OnSuccessK(getTrace(1), fn)
}

func OnSuccessIO[T1, T2 any](io IO[T1], fn func(T1) IO[T2]) IO[o.Option[T2]] {
	return _OnSuccessIO(getTrace(1), io, fn)
}
func OnSuccessIOK[T1, T2 any](fn func(T1) IO[T2]) func(IO[T1]) IO[o.Option[T2]] {
	return _OnSuccessIOK(getTrace(1), fn)
}

func Redeem[T any](io IO[T], fn func(error) T) IO[T] {
	return _Redeem(getTrace(1), io, fn)
}
func RedeemK[T any](fn func(error) T) func(IO[T]) IO[T] {
	return _RedeemK(getTrace(1), fn)
}

func RedeemSome[T any](io IO[T], fn func(error) o.Option[T]) IO[T] {
	return _RedeemSome(getTrace(1), io, fn)
}
func RedeemSomeK[T any](fn func(error) o.Option[T]) func(IO[T]) IO[T] {
	return _RedeemSomeK[T](getTrace(1), fn)
}

func RedeemIO[T any](io IO[T], fn func(error) IO[T]) IO[T] {
	return _RedeemIO(getTrace(1), io, fn)
}
func RedeemIOK[T any](fn func(error) IO[T]) func(IO[T]) IO[T] {
	return _RedeemIOK(getTrace(1), fn)
}

func RedeemSomeIO[T any](io IO[T], fn func(error) o.Option[IO[T]]) IO[T] {
	return _RedeemSomeIO(getTrace(1), io, fn)
}
func RedeemSomeIOK[T any](fn func(error) o.Option[IO[T]]) func(IO[T]) IO[T] {
	return _RedeemSomeIOK[T](getTrace(1), fn)
}

func MapError[T any](io IO[T], fn func(error) error) IO[T] {
	return _MapError(getTrace(1), io, fn)
}
func MapErrorK[T any](fn func(error) error) func(IO[T]) IO[T] {
	return _MapErrorK[T](getTrace(1), fn)
}

func When[T any](cond bool, io IO[T]) IO[o.Option[T]] {
	return _When(getTrace(1), cond, io)
}
func WhenK[T any](cond bool) func(IO[T]) IO[o.Option[T]] {
	return _WhenK[T](getTrace(1), cond)
}

func WhenIO[T any](cond IO[bool], io IO[T]) IO[o.Option[T]] {
	return _WhenIO(getTrace(1), cond, io)
}
func WhenIOK[T any](cond IO[bool]) func(IO[T]) IO[o.Option[T]] {
	return _WhenIOK[T](getTrace(1), cond)
}

func WhenM[T any](io IO[T]) func(bool) IO[o.Option[T]] {
	return _WhenM[T](getTrace(1), io)
}
func WhenIOM[T any](io IO[T]) func(IO[bool]) IO[o.Option[T]] {
	return _WhenIOM[T](getTrace(1), io)
}

func If[T any](cond bool, ifTrue IO[T], ifFalse IO[T]) IO[T] {
	return _If(getTrace(1), cond, ifTrue, ifFalse)
}
func IfIO[T any](cond IO[bool], ifTrue IO[T], ifFalse IO[T]) IO[T] {
	return _IfIO(getTrace(1), cond, ifTrue, ifFalse)
}
func IfK[T any](ifTrue IO[T], ifFalse IO[T]) func(bool) IO[T] {
	return _IfK(getTrace(1), ifTrue, ifFalse)
}
func IfIOK[T any](ifTrue IO[T], ifFalse IO[T]) func(IO[bool]) IO[T] {
	return _IfIOK(getTrace(1), ifTrue, ifFalse)
}

func Sequence[T any](ios []IO[T]) IO[[]T] {
	return _Traverse(getTrace(1), ios, f.Identity[IO[T]])
}
func ParSequence[T any](ios []IO[T], maxConcurrency o.Option[r.PosInt]) IO[[]T] {
	return _ParTraverse(getTrace(1), ios, f.Identity[IO[T]], maxConcurrency)
}

func Sequence_[T any](ios []IO[T]) VIO {
	trace := getTrace(1)
	return _As(trace, _Traverse(trace, ios, f.Identity[IO[T]]), v.Void{})
}
func ParSequence_[T any](ios []IO[T], maxConcurrency o.Option[r.PosInt]) VIO {
	trace := getTrace(1)
	return _As(trace, _ParTraverse(trace, ios, f.Identity[IO[T]], maxConcurrency), v.Void{})
}

func Traverse[T1, T2 any](ts []T1, liftFn func(T1) IO[T2]) IO[[]T2] {
	return _Traverse(getTrace(1), ts, liftFn)
}
func ParTraverse[T1, T2 any](ts []T1, liftFn func(T1) IO[T2], maxConcurrency o.Option[r.PosInt]) IO[[]T2] {
	return _ParTraverse(getTrace(1), ts, liftFn, maxConcurrency)
}

func Traverse_[T1, T2 any](ts []T1, liftFn func(T1) IO[T2]) VIO {
	return _Traverse(getTrace(1), ts, liftFn).Void()
}
func ParTraverse_[T1, T2 any](ts []T1, liftFn func(T1) IO[T2], maxConcurrency o.Option[r.PosInt]) VIO {
	trace := getTrace(1)
	return _As(trace, _ParTraverse(trace, ts, liftFn, maxConcurrency), v.Void{})
}

func Race[T any](ios []IO[T], maxConcurrency o.Option[r.PosInt]) IO[T] {
	return _Race(getTrace(1), ios, maxConcurrency)
}
func RacePair[T1, T2 any](left IO[T1], right IO[T2]) IO[e.Either[T1, T2]] {
	return _RacePair(getTrace(1), left, right)
}

func RunAll(ios ...VIO) VIO {
	return _RunAll(getTrace(1), ios...)
}

//Once Return an effect which will be executing the given effect at most once
func Once[T any](io IO[T]) IO[IO[T]] {
	return _Once(getTrace(1), io)
}

//Once_ Return an effect which will be executing the given effect at most once, discarding the value
func Once_[T any](io IO[T]) IO[VIO] {
	return _Once_(getTrace(1), io)
}

func NonInteruptible[T any](io IO[T]) IO[T] {
	return _Uninterruptible(getTrace(1), io)
}

func WithContext[T any](fn func(context.Context) IO[T]) IO[T] {
	return _WithContext[T](getTrace(1), fn)
}

func While(cond IO[bool], do VIO) VIO {
	return _While(getTrace(1), cond, do)
}

func TailRec[A, B any](curr A, do func(A) IO[e.Either[A, B]]) IO[B] {
	return _TailRec(getTrace(1), curr, do)
}

func Bracket[A, B any](acquire IO[A], use func(A) IO[B], release func(A) VIO) IO[B] {
	return _Bracket(getTrace(1), acquire, use, release)
}

func PartialUninterruptible[T any](io func(InterruptibilityContext) IO[T]) IO[T] {
	return _PartialUninterruptible(getTrace(1), io)
}
func RestoreInterruptibility[T any](context InterruptibilityContext, io IO[T]) IO[o.Option[T]] {
	return _RestoreInterruptibility(getTrace(1), context, io)
}

func Blocking[T any](io IO[T]) IO[T] {
	return _Blocking(getTrace(1), io)
}

func AsyncAll[T any](maxConcurrency o.Option[r.PosInt], asyncIos func(RunAsync AsyncAllRun) VIO, yield IO[T]) IO[T] {
	return _AsyncAll(getTrace(1), maxConcurrency, asyncIos, yield)
}

func Fork_[A any](_trace *Trace, io IO[A]) VIO {
	return _Fork(_trace, io).Void()
}

func Fork[A any](_trace *Trace, io IO[A]) IO[IO[A]] {
	return _Fork(_trace, io)
}

func UnsafeForkAndForget[A any](_trace *Trace, io IO[A]) VIO {
	return _UnsafeForkAndForget(_trace, io)
}

func Timed[A any](io IO[A]) IO[t.T2[time.Duration, A]] {
	return _Timed(getTrace(1), io)
}

func Void() VIO { return _Pure(v.Void{}) }

func AndThenK[A any](io IO[A]) func(any) IO[A] {
	return func(a any) IO[A] { return io }
}
