package io

import (
	"context"
	f "github.com/valentinHenry/giog/utils/functions"
	e "github.com/valentinHenry/giog/utils/monads/either"
	o "github.com/valentinHenry/giog/utils/monads/option"
	p "github.com/valentinHenry/giog/utils/pipes"
	t "github.com/valentinHenry/giog/utils/tuples"
	v "github.com/valentinHenry/giog/utils/void"
	r "github.com/valentinHenry/refined"
	"golang.org/x/sync/errgroup"
	"time"
)

func _Map[T1, T2 any](_trace *Trace, io IO[T1], mapFn func(T1) T2) IO[T2] {
	return _FlatMap(_trace, io, func(t1 T1) IO[T2] { return exitSuccess[T2](mapFn(t1)) })
}
func _MapK[T1, T2 any](_trace *Trace, mapFn func(T1) T2) func(IO[T1]) IO[T2] {
	return _FlatMapK(_trace, func(t1 T1) IO[T2] { return exitSuccess[T2](mapFn(t1)) })
}

func _MapError[T any](_trace *Trace, io IO[T], fn func(error) error) IO[T] {
	return &_IOFailure[T]{
		trace:    _trace,
		previous: io,
		onFailure: func(c Cause) IO[T] {
			mappedError := fn(c.Cause())
			return exitError[T](makeMappedCause(_trace, c, &mappedError))
		},
	}
}
func _MapErrorK[T any](_trace *Trace, fn func(error) error) func(IO[T]) IO[T] {
	return func(io IO[T]) IO[T] {
		return _MapError(_trace, io, fn)
	}
}

func _MapBoth[T1, T2 any](_trace *Trace, io IO[T1], onError func(error) error, onSuccess func(T1) T2) IO[T2] {
	return &_IOSuccessFailure[T1, T2]{
		trace:     _trace,
		previous:  io,
		onSuccess: func(t1 T1) IO[T2] { return _Pure(onSuccess(t1)) },
		onFailure: func(c Cause) IO[T2] {
			mappedError := onError(c.Cause())
			return exitError[T2](makeMappedCause(_trace, c, &mappedError))
		},
	}
}
func _MapBothK[T1, T2 any](_trace *Trace, onError func(error) error, onSuccess func(T1) T2) func(IO[T1]) IO[T2] {
	return func(io IO[T1]) IO[T2] {
		return _MapBoth(_trace, io, onError, onSuccess)
	}
}

func _FlatMap[T1, T2 any](_trace *Trace, io IO[T1], mapFn func(T1) IO[T2]) IO[T2] {
	return &_IOSuccess[T1, T2]{_trace, io, mapFn}
}
func _FlatMapK[T1, T2 any](_trace *Trace, mapFn func(T1) IO[T2]) func(IO[T1]) IO[T2] {
	return _IOSuccessK[T1, T2](_trace, mapFn)
}

func _FlatTap[T1, T2 any](_trace *Trace, io IO[T1], mapFn func(T1) IO[T2]) IO[T1] {
	return &_IOSuccess[T1, T1]{
		_trace,
		io,
		func(t1 T1) IO[T1] { return _As(_trace, mapFn(t1), t1) },
	}
}
func _FlatTapK[T1, T2 any](_trace *Trace, mapFn func(T1) IO[T2]) func(IO[T1]) IO[T1] {
	return _IOSuccessK[T1, T1](
		_trace,
		func(t1 T1) IO[T1] { return _As(_trace, mapFn(t1), t1) },
	)
}

func _As[T1, T2 any](_trace *Trace, io IO[T1], v T2) IO[T2] {
	return _Map(_trace, io, func(t T1) T2 { return v })
}
func _AsK[T1, T2 any](_trace *Trace, v T2) func(IO[T1]) IO[T2] {
	return _MapK(_trace, func(t T1) T2 { return v })
}

func _Absolve[T any](_trace *Trace, io IO[e.Either[error, T]]) IO[T] {
	return _FlatMap(_trace, io, _FromEitherK[T](_trace))
}

func _FromEither[T any](_trace *Trace, either e.Either[error, T]) IO[T] {
	return e.Fold(either, raiseK[T](_trace), exitSuccess[T])
}

func _FromEitherK[T any](_trace *Trace) func(e.Either[error, T]) IO[T] {
	return func(in e.Either[error, T]) IO[T] {
		return _FromEither(_trace, in)
	}
}

func _FromEitherDefer[T any](_trace *Trace, io func() IO[e.Either[error, T]]) IO[T] {
	return _FlatMap(
		_trace,
		_Defer(_trace, io),
		_FromEitherK[T](_trace),
	)
}

func _FromEitherDelay[T any](_trace *Trace, either func() e.Either[error, T]) IO[T] {
	return _FlatMap(
		_trace,
		_Delay(_trace, either),
		_FromEitherK[T](_trace),
	)
}

func _Pure[T any](a T) IO[T] {
	return exitSuccess(a)
}
func _PureK[T any]() func(T) IO[T] {
	return func(t T) IO[T] { return exitSuccess(t) }
}

func _Raise[T any](_trace *Trace, err error) IO[T] {
	return exitError[T](makeRaisedCause(_trace, &err))
}
func _RaiseK[T any](_trace *Trace) func(error) IO[T] {
	return raiseK[T](_trace)
}

func _Delay[T any](_trace *Trace, a func() T) IO[T] {
	return &_IOSync[T]{
		trace: _trace,
		eval:  func(_ context.Context) T { return a() },
	}
}
func _DelayK[T any](_trace *Trace) func(func() T) IO[T] {
	return func(a func() T) IO[T] {
		return _Delay(_trace, a)
	}
}

func _Defer[T any](_trace *Trace, io func() IO[T]) IO[T] {
	return _FlatMap(
		_trace,
		_Delay(_trace, io),
		f.Identity[IO[T]],
	)
}
func _DeferK[T any](_trace *Trace) func(func() IO[T]) IO[T] {
	return p.Pipe2K(
		_DelayK[IO[T]](_trace),
		_FlatMapK(_trace, f.Identity[IO[T]]),
	)
}

func _LiftV(_trace *Trace, fn func() error) VIO {
	return _Lift(_trace, func() (v.Void, error) { return v.Void{}, fn() })
}

func _LiftVK(_trace *Trace) func(func() error) VIO {
	return func(fn func() error) VIO {
		return _LiftV(_trace, fn)
	}
}

func _Lift[T any](_trace *Trace, v func() (T, error)) IO[T] {
	return _Defer(
		_trace,
		func() IO[T] {
			t, err := v()
			if err != nil {
				return _Raise[T](_trace, err)
			}
			return exitSuccess(t)
		},
	)
}
func _LiftK[T any](_trace *Trace) func(func() (T, error)) IO[T] {
	return func(v func() (T, error)) IO[T] {
		return _Lift[T](_trace, v)
	}
}

func _Fold[T1, T2 any](_trace *Trace, io IO[T1], onError func(error) T2, onSuccess func(T1) T2) IO[T2] {
	return &_IOSuccessFailure[T1, T2]{
		trace:     _trace,
		previous:  io,
		onSuccess: func(t T1) IO[T2] { return _Pure(onSuccess(t)) },
		onFailure: func(c Cause) IO[T2] { return _Pure(onError(c.Cause())) },
	}
}
func _FoldK[T1, T2 any](_trace *Trace, onError func(error) T2, onSuccess func(T1) T2) func(IO[T1]) IO[T2] {
	return func(io IO[T1]) IO[T2] {
		return _Fold[T1, T2](_trace, io, onError, onSuccess)
	}
}

func _FoldIO[T1, T2 any](_trace *Trace, io IO[T1], onError func(error) IO[T2], onSuccess func(T1) IO[T2]) IO[T2] {
	return &_IOSuccessFailure[T1, T2]{
		trace:     _trace,
		previous:  io,
		onSuccess: onSuccess,
		onFailure: func(c Cause) IO[T2] { return onError(c.Cause()) },
	}
}
func _FoldIOK[T1, T2 any](_trace *Trace, onError func(error) IO[T2], onSuccess func(T1) IO[T2]) func(IO[T1]) IO[T2] {
	return func(io IO[T1]) IO[T2] {
		return _FoldIO[T1, T2](_trace, io, onError, onSuccess)
	}
}

func _FoldCause[T1, T2 any](_trace *Trace, io IO[T1], onError func(Cause) T2, onSuccess func(T1) T2) IO[T2] {
	return &_IOSuccessFailure[T1, T2]{
		trace:     _trace,
		previous:  io,
		onSuccess: func(t T1) IO[T2] { return _Pure(onSuccess(t)) },
		onFailure: func(c Cause) IO[T2] { return _Pure(onError(c)) },
	}
}
func _FoldCauseK[T1, T2 any](_trace *Trace, onError func(Cause) T2, onSuccess func(T1) T2) func(IO[T1]) IO[T2] {
	return func(io IO[T1]) IO[T2] {
		return _FoldCause(_trace, io, onError, onSuccess)
	}
}

func _FoldIOCause[T1, T2 any](_trace *Trace, io IO[T1], onError func(Cause) IO[T2], onSuccess func(T1) IO[T2]) IO[T2] {
	return &_IOSuccessFailure[T1, T2]{
		trace:     _trace,
		previous:  io,
		onSuccess: onSuccess,
		onFailure: onError,
	}
}
func _FoldIOCauseK[T1, T2 any](_trace *Trace, onError func(Cause) IO[T2], onSuccess func(T1) IO[T2]) func(IO[T1]) IO[T2] {
	return func(io IO[T1]) IO[T2] {
		return _FoldIOCause(_trace, io, onError, onSuccess)
	}
}

func _OnSuccess[T1, T2 any](_trace *Trace, io IO[T1], fn func(T1) T2) IO[o.Option[T2]] {
	return _Fold(
		_trace,
		io,
		func(_ error) o.Option[T2] { return o.None[T2]{} },
		func(t1 T1) o.Option[T2] { return o.Some[T2]{fn(t1)} },
	)
}
func _OnSuccessK[T1, T2 any](_trace *Trace, fn func(T1) T2) func(IO[T1]) IO[o.Option[T2]] {
	return _FoldK(
		_trace,
		func(_ error) o.Option[T2] { return o.None[T2]{} },
		func(t1 T1) o.Option[T2] { return o.Some[T2]{fn(t1)} },
	)
}

func _OnSuccessIO[T1, T2 any](_trace *Trace, io IO[T1], fn func(T1) IO[T2]) IO[o.Option[T2]] {
	return _FoldIO(
		_trace,
		io,
		func(_ error) IO[o.Option[T2]] { return Pure[o.Option[T2]](o.None[T2]{}) },
		func(t1 T1) IO[o.Option[T2]] {
			return _Map(_trace, fn(t1), func(t2 T2) o.Option[T2] { return o.Some[T2]{t2} })
		},
	)
}
func _OnSuccessIOK[T1, T2 any](_trace *Trace, fn func(T1) IO[T2]) func(IO[T1]) IO[o.Option[T2]] {
	return _FoldIOK(
		_trace,
		func(_ error) IO[o.Option[T2]] { return Pure[o.Option[T2]](o.None[T2]{}) },
		func(t1 T1) IO[o.Option[T2]] {
			return _Map(_trace, fn(t1), func(t2 T2) o.Option[T2] { return o.Some[T2]{t2} })
		},
	)
}

func _Redeem[T any](_trace *Trace, io IO[T], fn func(error) T) IO[T] {
	return &_IOFailure[T]{
		trace:     _trace,
		previous:  io,
		onFailure: func(c Cause) IO[T] { return _Pure(fn(c.Cause())) },
	}
}
func _RedeemK[T any](_trace *Trace, fn func(error) T) func(IO[T]) IO[T] {
	return func(io IO[T]) IO[T] {
		return _Redeem(_trace, io, fn)
	}
}

func _RedeemSome[T any](_trace *Trace, io IO[T], fn func(error) o.Option[T]) IO[T] {
	return &_IOFailure[T]{
		trace:    _trace,
		previous: io,
		onFailure: func(c Cause) IO[T] {
			return o.Fold(
				fn(c.Cause()),
				exitError[T](c),
				_Pure[T],
			)
		},
	}
}
func _RedeemSomeK[T any](_trace *Trace, fn func(error) o.Option[T]) func(IO[T]) IO[T] {
	return func(io IO[T]) IO[T] {
		return _RedeemSome(_trace, io, fn)
	}
}

func _RedeemIO[T any](_trace *Trace, io IO[T], fn func(error) IO[T]) IO[T] {
	return &_IOFailure[T]{
		trace:     _trace,
		previous:  io,
		onFailure: func(c Cause) IO[T] { return fn(c.Cause()) },
	}
}
func _RedeemIOK[T any](_trace *Trace, fn func(error) IO[T]) func(IO[T]) IO[T] {
	return func(io IO[T]) IO[T] {
		return _RedeemIO(_trace, io, fn)
	}
}

func _RedeemSomeIO[T any](_trace *Trace, io IO[T], fn func(error) o.Option[IO[T]]) IO[T] {
	return &_IOFailure[T]{
		trace:    _trace,
		previous: io,
		onFailure: func(c Cause) IO[T] {
			return o.Fold(
				fn(c.Cause()),
				exitError[T](c),
				f.Identity[IO[T]],
			)
		},
	}
}
func _RedeemSomeIOK[T any](_trace *Trace, fn func(error) o.Option[IO[T]]) func(IO[T]) IO[T] {
	return func(io IO[T]) IO[T] {
		return _RedeemSomeIO(_trace, io, fn)
	}
}

func _When[T any](_trace *Trace, cond bool, io IO[T]) IO[o.Option[T]] {
	if cond {
		return _Map(_trace, io, func(t T) o.Option[T] { return o.Some[T]{t} })
	} else {
		return _Pure[o.Option[T]](o.None[T]{})
	}
}
func _WhenK[T any](_trace *Trace, cond bool) func(IO[T]) IO[o.Option[T]] {
	return func(io IO[T]) IO[o.Option[T]] {
		return _When(_trace, cond, io)
	}
}
func _WhenM[T any](_trace *Trace, io IO[T]) func(bool) IO[o.Option[T]] {
	return func(cond bool) IO[o.Option[T]] {
		return _When(_trace, cond, io)
	}
}

func _WhenIO[T any](_trace *Trace, cond IO[bool], io IO[T]) IO[o.Option[T]] {
	return _FlatMap(_trace, cond, func(cond bool) IO[o.Option[T]] {
		if cond {
			return _Map(_trace, io, func(t T) o.Option[T] { return o.Some[T]{t} })
		} else {
			return _Pure[o.Option[T]](o.None[T]{})
		}
	})
}
func _WhenIOK[T any](_trace *Trace, cond IO[bool]) func(IO[T]) IO[o.Option[T]] {
	return func(io IO[T]) IO[o.Option[T]] {
		return _WhenIO(_trace, cond, io)
	}
}
func _WhenIOM[T any](_trace *Trace, io IO[T]) func(IO[bool]) IO[o.Option[T]] {
	return func(cond IO[bool]) IO[o.Option[T]] {
		return _WhenIO(_trace, cond, io)
	}
}

func _If[T any](_trace *Trace, cond bool, ifTrue IO[T], ifFalse IO[T]) IO[T] {
	return _IfIO(_trace, _Pure(cond), ifTrue, ifFalse)
}
func _IfK[T any](_trace *Trace, ifTrue IO[T], ifFalse IO[T]) func(bool) IO[T] {
	return func(cond bool) IO[T] {
		return _If(_trace, cond, ifTrue, ifFalse)
	}
}

func _IfIO[T any](_trace *Trace, cond IO[bool], ifTrue IO[T], ifFalse IO[T]) IO[T] {
	return _FlatMap(
		_trace,
		cond,
		func(cond bool) IO[T] {
			if cond {
				return ifTrue
			} else {
				return ifFalse
			}
		},
	)
}
func _IfIOK[T any](_trace *Trace, ifTrue IO[T], ifFalse IO[T]) func(IO[bool]) IO[T] {
	return func(cond IO[bool]) IO[T] {
		return _IfIO(_trace, cond, ifTrue, ifFalse)
	}
}

func _Traverse[T1, T2 any](_trace *Trace, ts []T1, liftFn func(T1) IO[T2]) IO[[]T2] {
	return _TraverseZip(_trace, ts, func(_ int, t1 T1) IO[T2] { return liftFn(t1) })
}
func _TraverseZip[T1, T2 any](_trace *Trace, ts []T1, liftFn func(int, T1) IO[T2]) IO[[]T2] {
	// Recursive version of _Traverse, although it works, for more efficiency _While is used
	//func _Traverse[T1, T2 any](_trace *Trace, ts []T1, liftFn func(T1) IO[T2]) IO[[]T2] {
	//	resTs := make([]T2, len(ts))
	//	return traverse(_trace, 0, resTs, liftFn, ts)
	//}
	//func traverse[T1, T2 any](_trace *Trace, i int, res []T2, liftFn func(T1) IO[T2], remaining []T1) IO[[]T2] {
	//	if len(remaining) == 0 {
	//		return _Pure(res)
	//	}
	//	return _FlatMap(
	//		_trace,
	//		liftFn(remaining[0]),
	//		func(t2 T2) IO[[]T2] {
	//			res[i] = t2
	//			return traverse(_trace, i+1, res, liftFn, remaining[1:])
	//		},
	//	)
	//}
	l := len(ts)
	curr := 0
	res := make([]T2, l)

	loop :=
		_While(
			_trace,
			_Delay(_trace, func() bool { return curr < l }),
			_Map(
				_trace,
				Defer(func() IO[T2] { return liftFn(curr, ts[curr]) }),
				func(r T2) v.Void {
					res[curr] = r
					curr = curr + 1
					return v.Void{}
				},
			),
		)

	return _AndThen2(
		_trace,
		loop,
		_Delay(_trace, func() []T2 { return res }),
	)
}

func _ParTraverse[T1, T2 any](_trace *Trace, ts []T1, liftFn func(T1) IO[T2], maxConcurrency o.Option[r.PosInt]) IO[[]T2] {
	setValue := func(ref Ref[[]T2], idx int, RunAsync AsyncAllRun) VIO {
		return _FlatMap(
			_trace,
			liftFn(ts[idx]),
			func(res T2) VIO {
				return RunAsync(ref.Update(func(arr []T2) []T2 { arr[idx] = res; return arr }))
			},
		)
	}

	curr := 0
	l := len(ts)

	runAll := func(ref Ref[[]T2], RunAsync AsyncAllRun) VIO {
		return _While(
			_trace,
			_Delay(_trace, func() bool { return curr < l }),
			_AndThen2(
				_trace,
				_Defer(_trace, func() VIO { return setValue(ref, curr, RunAsync) }),
				_Delay(_trace, func() v.Void {
					curr++
					return v.Void{}
				}),
			),
		)
	}

	return _FlatMap(
		_trace,
		MakeRef(make([]T2, len(ts))),
		func(ref Ref[[]T2]) IO[[]T2] {
			return _AsyncAll(
				_trace,
				maxConcurrency,
				func(ctx AsyncAllRun) VIO { return runAll(ref, ctx) },
				ref.Get(),
			)
		},
	)
}

func _Race[T any](_trace *Trace, ios []IO[T], maxConcurrency o.Option[r.PosInt]) IO[T] {
	raiseOnSuccess := func(io IO[T]) VIO {
		return _FlatMap(
			_trace,
			_InterruptFast(_trace, io),
			func(t T) VIO { return _Raise[v.Void](_trace, raceValue[T]{t}) },
		)
	}

	runAllAsync := func(ctx AsyncAllRun) VIO {
		return _Traverse(
			_trace,
			ios,
			func(io IO[T]) VIO { return ctx(raiseOnSuccess(io)) },
		).Void()
	}

	raceAll := _AsyncAll(
		_trace,
		maxConcurrency,
		runAllAsync,
		_Delay[T](_trace, func() T { panic("impossible") }), // Safe here since the value is contained on the error-side of the IO, it will never happen
	)

	return _RedeemSome(
		_trace,
		raceAll,
		func(err error) o.Option[T] {
			yield, ok := err.(raceValue[T])
			return o.When[T](ok, yield.value)
		},
	)
}

type raceValue[T any] struct {
	value T
}

func (r raceValue[T]) Error() string { return "" }

func _RacePair[T1, T2 any](_trace *Trace, v1 IO[T1], v2 IO[T2]) IO[e.Either[T1, T2]] {
	return _Race(
		_trace,
		[]IO[e.Either[T1, T2]]{
			_Map(_trace, v1, func(v T1) e.Either[T1, T2] { return e.Left[T1, T2]{v} }),
			_Map(_trace, v2, func(v T2) e.Either[T1, T2] { return e.Right[T1, T2]{v} }),
		},
		o.None[r.PosInt]{},
	)
}

func _RunAll(_trace *Trace, ios ...VIO) VIO {
	return _AsyncAll(
		_trace,
		o.None[r.PosInt]{},
		func(RunAsync AsyncAllRun) VIO {
			return _Traverse(_trace, ios, func(io VIO) VIO { return RunAsync(io) }).Void()
		},
		_Pure(v.Void{}),
	)
}

func _AsyncAll[T any](_trace *Trace, maxConcurrency o.Option[r.PosInt], asyncIos func(RunAsync AsyncAllRun) VIO, yield IO[T]) IO[T] {
	return _WithContext(_trace, func(ctx context.Context) IO[T] {
		g, gctx := errgroup.WithContext(ctx)
		g.SetLimit(o.Fold(maxConcurrency, -1, func(c r.PosInt) int { return c.Value() }))

		wait := _FlatMap(
			_trace,
			_InterruptFast(_trace, _Delay(_trace, g.Wait)),
			func(e error) VIO {
				if e == nil {
					return _Pure(v.Void{})
				}

				err, ok := e.(causeError)
				if ok {
					return exitError[v.Void](err.cause)
				}
				return _Raise[v.Void](_trace, err)
			},
		)

		var asyncCtx AsyncAllRun = func(io VIO) VIO {
			return &_IOAsync[v.Void]{
				trace:       _trace,
				ctx:         gctx,
				io:          io,
				forgettable: false,
				runAsync: func(toRun func() (v.Void, Cause)) {
					g.Go(func() error {
						_, err := toRun()
						if err != nil {
							return causeError{err}
						}
						return nil
					})
				},
			}
		}

		return _AndThen3(
			_trace,
			asyncIos(asyncCtx),
			wait,
			yield,
		)
	})
}

type causeError struct {
	cause Cause
}

func (e causeError) Error() string { return e.cause.Cause().Error() }

type AsyncAllRun func(VIO) VIO

func _Go[A any](_trace *Trace, io IO[A]) IO[IO[A]] {
	return _WithContext(
		_trace,
		func(ctx context.Context) IO[IO[A]] {
			resChan := make(chan e.Either[Cause, A])

			var runAsync VIO = &_IOAsync[A]{
				trace:       _trace,
				ctx:         ctx,
				io:          io,
				forgettable: false,
				runAsync: func(toRun func() (A, Cause)) {
					go func() {
						a, err := toRun()
						if err != nil {
							resChan <- e.Left[Cause, A]{err}
						} else {
							resChan <- e.Right[Cause, A]{a}
						}
					}()
				},
			}

			var waitForResult IO[A] = _Defer(_trace,
				func() IO[A] {
					select {
					case <-ctx.Done():
						return exitError[A](makeCancellationCause())
					case res := <-resChan:
						return e.Fold(res,
							func(c Cause) IO[A] { return exitError[A](c) },
							func(v A) IO[A] { return _Pure(v) },
						)
					}
				},
			)

			return _As(_trace, runAsync, waitForResult)
		},
	)
}

func _UnsafeGo_[A any](_trace *Trace, io IO[A]) VIO {
	return _WithContext(
		_trace,
		func(ctx context.Context) VIO {
			return &_IOAsync[A]{
				trace:       _trace,
				ctx:         ctx,
				io:          io,
				forgettable: true,
				runAsync: func(toRun func() (A, Cause)) {
					go func() {
						_, _ = toRun()
					}()
				},
			}
		},
	)
}

func _Pair2[T1, T2 any](_trace *Trace, v1 IO[T1], v2 IO[T2]) IO[t.T2[T1, T2]] {
	return _FlatMap(_trace, v1, func(v1 T1) IO[t.T2[T1, T2]] {
		return _Map(_trace, v2, func(v2 T2) t.T2[T1, T2] {
			return t.Of2(v1, v2)
		})
	})
}

func _Once[T any](_trace *Trace, io IO[T]) IO[IO[T]] {
	condRef := MakeRef(false)
	resVar := MakeDeferred[T]()

	var getAndUpdate f.Fn2[Ref[bool], Deferred[T], IO[T]] = func(condRef Ref[bool], resVar Deferred[T]) IO[T] {
		return _IfIO(
			_trace,
			condRef.GetAndSet(true),
			/* ifTrue  */ resVar.Get(), // Blocks in case the execution of the effect has not completed yet
			/* ifFalse */ _FlatTap(_trace, io, resVar.Complete),
		)
	}

	return _Map(
		_trace,
		_Pair2(_trace, condRef, resVar),
		getAndUpdate.Tupled,
	)
}
func _Once_[T any](_trace *Trace, io IO[T]) IO[VIO] {
	return _Map(_trace,
		MakeRef(true),
		func(ref Ref[bool]) VIO {
			return _WhenIO(_trace, ref.GetAndSet(false), io).Void()
		},
	)
}

func _Uncancelable[T any](_trace *Trace, io IO[T]) IO[T] {
	return &_IOUniverseSwitch[T]{
		trace:         _trace,
		get:           func(universe *Universe) *Universe { return universe.CloneWithContext(universe.Uninterruptible) },
		withUniverses: func(_ *Universe, _ *Universe) IO[T] { return io },
		release:       func(_ *Universe) {},
	}
}

func _PartialUncancelable[T any](_trace *Trace, io func(CancelabilityContext) IO[T]) IO[T] {
	return &_IOUniverseSwitch[T]{
		trace: _trace,
		get: func(universe *Universe) *Universe {
			return universe.CloneWithContext(universe.Uninterruptible)
		},
		withUniverses: func(old *Universe, new *Universe) IO[T] {
			return io(CancelabilityContext{context: old.Context})
		},
		release: func(_ *Universe) {},
	}
}

func _RestoreCancelability[T any](_trace *Trace, context CancelabilityContext, io IO[T]) IO[T] {
	cancellable := &_IOUniverseSwitch[T]{
		trace:         _trace,
		get:           func(universe *Universe) *Universe { return universe.CloneWithContext(context.context) },
		withUniverses: func(_ *Universe, _ *Universe) IO[T] { return io },
		release:       func(_ *Universe) {},
	}

	return cancellable
}

func _OnCancelled[T any](_trace *Trace, io IO[T], ifCancelled IO[T]) IO[T] {
	return _PartialUncancelable(_trace, func(ctx CancelabilityContext) IO[T] {
		return &_IOOnCancel[T]{
			trace:    _trace,
			previous: _RestoreCancelability(_trace, ctx, io),
			onCancel: ifCancelled,
		}
	})
}

type CancelabilityContext struct {
	context context.Context
}

func _WithContext[T any](_trace *Trace, fn func(context.Context) IO[T]) IO[T] {
	return _FlatMap(
		_trace,
		(IO[IO[T]])(&_IOSync[IO[T]]{trace: _trace, eval: fn}),
		f.Identity[IO[T]],
	)
}

func _While(_trace *Trace, cond IO[bool], do VIO) VIO {
	return &_IOWhileLoop{trace: _trace, cond: cond, do: do}
}

func _TailRec[A, B any](_trace *Trace, curr A, do func(A) IO[e.Either[A, B]]) IO[B] {
	return _FlatMap(
		_trace,
		do(curr),
		func(res e.Either[A, B]) IO[B] {
			return e.Fold(
				res,
				func(curr A) IO[B] { return _TailRec(_trace, curr, do) },
				_Pure[B],
			)
		},
	)
}

func _Bracket[A, B any](_trace *Trace, acquire IO[A], use func(A) IO[B], release func(A) VIO) IO[B] {
	return _PartialUncancelable(
		_trace,
		func(restorer CancelabilityContext) IO[B] {
			runAndRelease := func(a A) IO[B] {
				var runAndReleaseUninterruptibe IO[B] = _FoldIOCause(
					_trace,
					_RestoreCancelability(_trace, restorer, use(a)),
					func(cause Cause) IO[B] { return _AndThen2(_trace, release(a), exitError[B](cause)) },
					func(res B) IO[B] { return _As(_trace, release(a), res) },
				)

				return _OnCancelled(
					_trace,
					runAndReleaseUninterruptibe,
					_AndThen2(_trace, release(a), exitError[B](makeCancellationCause())), // Rethrowing cancellation
				)
			}

			return _FlatMap(
				_trace,
				acquire,
				runAndRelease,
			)
		},
	)
}

func _InterruptFast[A any](_trace *Trace, io IO[A]) IO[A] {
	return &_IOInterruptFast[A]{
		trace: _trace,
		io:    io,
	}
}

func _Blocking[A any](_trace *Trace, io IO[A]) IO[A] {
	return _InterruptFast[A](_trace, io)
}

func cancelledOpt[A any](_trace *Trace, io IO[A]) IO[o.Option[A]] {
	return &_IOOnCancel[o.Option[A]]{
		trace:    _trace,
		previous: _Map(_trace, io, func(a A) o.Option[A] { return o.Some[A]{a} }),
		onCancel: _Pure[o.Option[A]](o.None[A]{}),
	}
}

func _Timed[A any](_trace *Trace, io IO[A]) IO[t.T2[time.Duration, A]] {
	runWithTime := _Accumulate3(
		_trace,
		_Delay(_trace, func() time.Time { return time.Now() }),
		io,
		_Delay(_trace, func() time.Time { return time.Now() }),
	)

	var withDuration f.Fn3[time.Time, A, time.Time, t.T2[time.Duration, A]] = func(start time.Time, res A, end time.Time) t.T2[time.Duration, A] {
		return t.Of2(end.Sub(start), res)
	}

	return _Map(_trace, runWithTime, withDuration.Tupled)
}

func _FromGo[A any](_trace *Trace, fn func(ctx context.Context, callback func(A, error))) IO[A] {
	return _WithContext(_trace, func(ctx context.Context) IO[A] {
		resChan := make(chan e.Either[error, A])

		callback := func(res A, error error) {
			resChan <- e.FromResult(error, res)
		}

		fn(ctx, callback)

		select {
		case <-ctx.Done():
			return exitError[A](makeCancellationCause())
		case r := <-resChan:
			return _FromEither(_trace, r)
		}
	})
}
