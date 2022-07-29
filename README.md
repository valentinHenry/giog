# GioG - IO Monad for Go

GioG is a library implementing an IO Monad for Go. It is shipped with other useful
monads like `Either` and `Option`.

## Table of contents
- [IO](#io)
  - [What is an IO ?](#what-is-an-io-)
  - [Example](#example)
  - [Functions](#functions)
    - [Mapping functions](#mapping-functions)
    - [Transformation functions](#transformation-functions)
    - [Maker functions](#maker-functions)
    - [Error Handling](#error-handling)
    - [Conditional functions](#conditional-functions)
    - [Slice functions](#slice-functions)
    - [Contextual functions](#contextual-functions)
    - [Loops](#loops)
    - [Others](#others)
  - [Asynchronicity](#asynchronicity)
  - [Composition Functions](#composition-functions)
  - [Evaluates the IO](#evaluates-the-io)
- [RIO](#rio)
  - [Basics](#basics)
  - [Composition Functions](#composition-functions)
- [DataTypes](#datatypes)
  - [CountDownLatch](#countdownlatch)
    - [Example](#example)
  - [CyclicBarrier](#cyclicbarrier)
    - [Example](#example)
  - [Deferred](#deferred)
    - [Example](#example)
  - [Queue](#queue)
    - [Example](#example)
  - [Ref](#ref)
  - [Semaphore](#semaphore)
- [Monads](#monads)
  - [Either](#either)
  - [Option](#option)
- [Helpers](#helpers)
  - [Functions](#functions)
  - [Pipes](#pipes)
  - [Tuples](#tuples)## IO
```go
type IO[A any] interface {
  UnsafeRun() (A, error)
  Void() IO[v.Void]
}
```
### What is an IO ?
GioG allow you to write code in a more functional way by encapsulating side effects,
giving a better error handling and propagation DSL. It is also shipped with asynchronous
functions allowing to create functional concurrent code.

### Example

Below is an example to one of the most basic program using IO: Hello World! 
```go
// example.go
package main

import (
  "fmt"
  "github.com/valentinHenry/giog/io/console"
  "github.com/valentinHenry/giog/io/io"
  p "github.com/valentinHenry/giog/utils/pipes"
)

func main() {
  mainLogic :=
    io.AndThen2(
      console.Println("Hey, what's your name ?"),
      io.FlatMap(
        console.ReadLn(),
        func(name string) io.VIO { return console.Printf("Hello %s!", name).Void() },
      ),
    )
  // Or using a Pipe combinator
  mainLogic = p.Pipe2(
    console.Println("Hey, what's your name ?"),
    io.FlatMapK(func(_ int) io.IO[string] { return console.ReadLn() }),
    io.FlatMapK(func(name string) io.VIO { return console.Printf("Hello %s!\n", name).Void() }),
  )

  _, err := io.EvalSync(mainLogic)
  if err != nil {
    fmt.Println(err.Error())
  }
}
```
```shell
$> go build example.go
$> ./example
Hey, what's your name?
Paul
Hello Paul!
```

### Functions 
Most of the functions have a K function. This K function is returning another
function taking an IO as parameter and doing the same as the non-K function.

All the functions below are declared in the [functions.go](https://github.com/valentinHenry/giog/blob/master/io/io/functions.go) file.

#### Mapping functions
The function below allows you to map over an effect.
```go
func Map[T1, T2 any](io IO[T1], mapFn func(T1) T2) IO[T2]
func MapK[T1, T2 any](mapFn func(T1) T2) func(IO[T1]) IO[T2]
func FlatMap[T1, T2 any](io IO[T1], mapFn func(T1) IO[T2]) IO[T2]
func FlatMapK[T1, T2 any](mapFn func(T1) IO[T2]) func(IO[T1]) IO[T2]
func FlatTap[T1, T2 any](io IO[T1], mapFn func(T1) IO[T2]) IO[T1]
func FlatTapK[T1, T2 any](mapFn func(T1) IO[T2]) func(IO[T1]) IO[T1]
func MapBoth[T1, T2 any](io IO[T1], onError func(error) error, onSuccess func(T1) T2) IO[T2]
func MapBothK[T1, T2 any](onError func(error) error, onSuccess func(T1) T2) func(IO[T1]) IO[T2]
func Fold[T1, T2 any](io IO[T1], onError func(error) T2, onSuccess func(T1) T2) IO[T2]
func FoldK[T1, T2 any](onError func(error) T2, onSuccess func(T1) T2) func(IO[T1]) IO[T2]
func FoldIO[T1, T2 any](io IO[T1], onError func(error) IO[T2], onSuccess func(T1) IO[T2]) IO[T2]
func FoldIOK[T1, T2 any](onError func(error) IO[T2], onSuccess func(T1) IO[T2]) func(IO[T1]) IO[T2]
func FoldCause[T1, T2 any](io IO[T1], onError func(Cause) T2, onSuccess func(T1) T2) IO[T2]
func FoldCauseK[T1, T2 any](onError func(Cause) T2, onSuccess func(T1) T2) func(IO[T1]) IO[T2]
func FoldIOCause[T1, T2 any](io IO[T1], onError func(Cause) IO[T2], onSuccess func(T1) IO[T2]) IO[T2]
func FoldIOCauseK[T1, T2 any](onError func(Cause) IO[T2], onSuccess func(T1) IO[T2]) func(IO[T1]) IO[T2]
func OnSuccess[T1, T2 any](io IO[T1], fn func(T1) T2) IO[o.Option[T2]]
func OnSuccessK[T1, T2 any](fn func(T1) T2) func(IO[T1]) IO[o.Option[T2]]
func OnSuccessIO[T1, T2 any](io IO[T1], fn func(T1) IO[T2]) IO[o.Option[T2]]
func OnSuccessIOK[T1, T2 any](fn func(T1) IO[T2]) func(IO[T1]) IO[o.Option[T2]]
```

#### Transformation functions
```go
func Flatten[T any](io IO[IO[T]]) IO[T]
func FlattenK[T any]() func(io IO[IO[T]]) IO[T]
func As[T1, T2 any](io IO[T1], v T2) IO[T2]
func AsK[T1, T2 any](v T2) func(IO[T1]) IO[T2]
func Absolve[T any](io IO[e.Either[error, T]]) IO[T]
```

#### Maker functions
```go
// Succeeses
func Pure[T any](v T) IO[T]
func PureK[T any]() func(T) IO[T]
func Delay[T any](v func() T) IO[T]
func DelayK[T any]() func(func() T) IO[T]
func DeferK[T any]() func(func() IO[T]) IO[T]
// Errors
func Raise[T any](err error) IO[T]
func RaiseK[T any]() func(error) IO[T]
// Succeeses or errors
func Lift[T any](v func() (T, error)) IO[T]
func LiftK[T any]() func(func() (T, error)) IO[T]
func LiftV(v func() error) VIO
func LiftVK() func(func() error) VIO
func FromEither[T any](either e.Either[error, T]) IO[T] 
func FromEitherK[T any]() func(e.Either[error, T]) IO[T] 
func FromEitherDelay[T any](either func() e.Either[error, T]) IO[T] 
func FromEitherDefer[T any](either func() IO[e.Either[error, T]]) IO[T]
// From io creation
func Defer[T any](io func() IO[T]) IO[T]
func WithContext[T any](fn func(context.Context) IO[T]) IO[T]
// Other
func Void() VIO
func FromGo[A any](fn func(ctx context.Context, callback func(A, error))) IO[A]
```

#### Error Handling
In addition to `Fold` functions which can be found [There](#mapping-functions). Other functions are 
available.
```go
func Redeem[T any](io IO[T], fn func(error) T) IO[T]
func RedeemK[T any](fn func(error) T) func(IO[T]) IO[T]
func RedeemSome[T any](io IO[T], fn func(error) o.Option[T]) IO[T]
func RedeemSomeK[T any](fn func(error) o.Option[T]) func(IO[T]) IO[T]
func RedeemIO[T any](io IO[T], fn func(error) IO[T]) IO[T]
func RedeemIOK[T any](fn func(error) IO[T]) func(IO[T]) IO[T]
func RedeemSomeIO[T any](io IO[T], fn func(error) o.Option[IO[T]]) IO[T]
func RedeemSomeIOK[T any](fn func(error) o.Option[IO[T]]) func(IO[T]) IO[T]
func MapError[T any](io IO[T], fn func(error) error) IO[T]
func MapErrorK[T any](fn func(error) error) func(IO[T]) IO[T]
func OnCancelled[T any](io IO[T], ifCancelled IO[T]) IO[T]
```

#### Conditional functions
```go
func When[T any](cond bool, io IO[T]) IO[o.Option[T]]
func WhenK[T any](cond bool) func(IO[T]) IO[o.Option[T]]
func WhenIO[T any](cond IO[bool], io IO[T]) IO[o.Option[T]]
func WhenIOK[T any](cond IO[bool]) func(IO[T]) IO[o.Option[T]]
func WhenM[T any](io IO[T]) func(bool) IO[o.Option[T]]
func WhenIOM[T any](io IO[T]) func(IO[bool]) IO[o.Option[T]]
func If[T any](cond bool, ifTrue IO[T], ifFalse IO[T]) IO[T]
func IfIO[T any](cond IO[bool], ifTrue IO[T], ifFalse IO[T]) IO[T]
func IfK[T any](ifTrue IO[T], ifFalse IO[T]) func(bool) IO[T]
func IfIOK[T any](ifTrue IO[T], ifFalse IO[T]) func(IO[bool]) IO[T]
```

#### Slice functions
```go
func Sequence[T any](ios []IO[T]) IO[[]T]
func ParSequence[T any](ios []IO[T], maxConcurrency o.Option[r.PosInt]) IO[[]T]
func Sequence_[T any](ios []IO[T]) VIO
func ParSequence_[T any](ios []IO[T], maxConcurrency o.Option[r.PosInt]) VIO
func Traverse[T1, T2 any](ts []T1, liftFn func(T1) IO[T2]) IO[[]T2]
func ParTraverse[T1, T2 any](ts []T1, liftFn func(T1) IO[T2], maxConcurrency o.Option[r.PosInt]) IO[[]T2]
func Traverse_[T1, T2 any](ts []T1, liftFn func(T1) IO[T2]) VIO
func ParTraverse_[T1, T2 any](ts []T1, liftFn func(T1) IO[T2], maxConcurrency o.Option[r.PosInt]) VIO
```

#### Contextual functions
```go
func Once[T any](io IO[T]) IO[IO[T]]
func Once_[T any](io IO[T]) IO[VIO]
func Uncancelable[T any](io IO[T]) IO[T]
func PartialUncancelable[T any](io func(CancelabilityContext) IO[T]) IO[T]
func RestoreCancelability[T any](context CancelabilityContext, io IO[T]) IO[T]
func Blocking[T any](io IO[T]) IO[T]
```

#### Loops
```go
func While(cond IO[bool], do VIO) VIO
func TailRec[A, B any](curr A, do func(A) IO[e.Either[A, B]]) IO[B]
```

#### Others
```go
func Bracket[A, B any](acquire IO[A], use func(A) IO[B], release func(A) VIO) IO[B]
func Timed[A any](io IO[A]) IO[t.T2[time.Duration, A]]
func AndThenK[In, A any](io IO[A]) func(In) IO[A]
func AndThenTapK[In, A any](io IO[A]) func(IO[In]) IO[In]
```

### Asynchronicity
```go
func AsyncAll[T any](maxConcurrency o.Option[r.PosInt], asyncIos func(RunAsync AsyncAllRun) VIO, yield IO[T]) IO[T]
func Go[A any](io IO[A]) IO[IO[A]]
func Go_[A any](io IO[A]) VIO
func UnsafeGo_[A any](io IO[A]) VIO
```

### Composition Functions
Chaining of IOs can be made with `FlatMap`, `FlatTap`, `Sequence`... However some helpers are provided
to ease the use of these functions and avoid working into infinitely nested functions.

The functions below are available from 1 to 30 parameters.

```go
func AndThen3[T1, T2, T3 any](v1 IO[T1], v2 IO[T2], v3 IO[T3]) IO[T3]
func Accumulate3[T1, T2, T3 any](v1 IO[T1], v2 IO[T2], v3 IO[T3]) IO[tuples.T3[T1, T2, T3]]
func Bind3[T1, T2, T3 any](v1 IO[T1], v2 func(T1) IO[T2], v3 func(T1, T2) IO[T3]) IO[tuples.T3[T1, T2, T3]]
func BindY3[T1, T2, T3 any](v1 IO[T1], v2 func(T1) IO[T2], v3 func(T1, T2) IO[T3]) IO[T3]
```

An IO can also be composed using the `Pipe` composition function available in the `pipes` package.

### Evaluates the IO
```go
func EvalAsync[A any](io IO[A]) (A, error)
func EvalSync[A any](io IO[A]) (A, error)
```

Two functions are available for IO evaluation:
- `EvalSync` which will wait on all goroutines to end until returning the value..
- `EvanAsync` runs the IO without waiting for goroutines to end.

## RIO
### Basics
### Composition Functions
## DataTypes
### CountDownLatch
```go
type CountDownLatch interface {
  Release() VIO
  Await() VIO
}

func MakeCountDownLatch(nb uint) IO[CountDownLatch]
```
A `CountDownLatch` is an interface that semantically blocks any goroutines
which waits on it. These are blocked until all defined latches are
released.

After all latches are released, the latch count is not reset
(unlike what CyclicBarrier does with waiters). Thus `Await()` will not be
blocking.

#### Example
```go
// example.go
package main

import (
  "fmt"
  c "github.com/valentinHenry/giog/io/console"
  "github.com/valentinHenry/giog/io/io"
  p "github.com/valentinHenry/giog/utils/pipes"
)

func main() {
  makeWaiter := func(latch io.CountDownLatch, nb int) io.VIO {
    return io.AndThen3(
      c.Printf("[%d] Waiting\n", nb),
      latch.Await(),
      c.Printf("[%d] I'm free!\n", nb),
    ).Void()
  }

  releaseLatch := func(latch io.CountDownLatch) io.VIO {
    return io.AndThen2(
      c.Println("Releasing a latch"),
      latch.Release(),
    )
  }

  mainLogic := p.Pipe8(
    io.MakeCountDownLatch(3),
    io.FlatTapK(func(latch io.CountDownLatch) io.VIO { return io.Go_(makeWaiter(latch, 1)) }),
    io.FlatTapK(func(latch io.CountDownLatch) io.VIO { return io.Go_(makeWaiter(latch, 2)) }),
    io.FlatTapK(releaseLatch),
    io.FlatTapK(func(latch io.CountDownLatch) io.VIO { return io.Go_(makeWaiter(latch, 3)) }),
    io.FlatTapK(releaseLatch),
    io.FlatTapK(func(latch io.CountDownLatch) io.VIO { return io.Go_(makeWaiter(latch, 4)) }),
    io.FlatTapK(releaseLatch),
    io.FlatTapK(func(latch io.CountDownLatch) io.VIO { return io.Go_(makeWaiter(latch, 4)) }),
  )

  _, err := io.EvalSync(mainLogic)
  if err != nil {
    fmt.Println(err.Error())
  }
}
```
```shell
$> go run example.go
Releasing a latch
[2] Waiting
[1] Waiting
Releasing a latch
[3] Waiting
Releasing a latch
[4] Waiting
[4] I'm free!
[1] I'm free!
[3] I'm free!
[5] Waiting
[5] I'm free!
[2] I'm free!
```

---
### CyclicBarrier
```go
type CyclicBarrier interface {
  Await() VIO
}

func MakeCyclicBarrier(parties uint) IO[CyclicBarrier]
```
A CyclicBarrier is an interface to a synchronizer which allows goroutines
to wait for each-others at a fixed point.

#### Example
```go
// example.go
package main

import (
  "fmt"
  c "github.com/valentinHenry/giog/io/console"
  "github.com/valentinHenry/giog/io/io"
  p "github.com/valentinHenry/giog/utils/pipes"
)

func main() {
  makeWaiter := func(barrier io.CyclicBarrier, nb int) io.VIO {
    return io.AndThen3(
      c.Printf("[%d] Waiting\n", nb),
      barrier.Await(),
      c.Printf("[%d] I'm free!\n", nb),
    ).Void()
  }

  mainLogic := p.Pipe6(
    io.MakeCyclicBarrier(3),
    io.FlatTapK(func(b io.CyclicBarrier) io.VIO { return io.Go_(makeWaiter(b, 1)) }),
    io.FlatTapK(func(b io.CyclicBarrier) io.VIO { return io.Go_(makeWaiter(b, 2)) }),
    io.FlatTapK(func(b io.CyclicBarrier) io.VIO { return io.Go_(makeWaiter(b, 3)) }),
    io.FlatTapK(func(b io.CyclicBarrier) io.VIO { return io.Go_(makeWaiter(b, 4)) }),
    io.FlatTapK(func(b io.CyclicBarrier) io.VIO { return io.Go_(makeWaiter(b, 5)) }),
    io.FlatTapK(func(b io.CyclicBarrier) io.VIO { return io.Go_(makeWaiter(b, 6)) }),
  )

  _, err := io.EvalSync(mainLogic)
  if err != nil {
    fmt.Println(err.Error())
  }
}
```
```shell
$> go run example.go
[1] Waiting
[6] Waiting
[3] Waiting
[3] I'm free!
[6] I'm free!
[1] I'm free!
[2] Waiting
[5] Waiting
[4] Waiting
[4] I'm free!
[5] I'm free!
[2] I'm free!
```

---
### Deferred
```go
type Deferred[A any] interface {
  Get() IO[A]
  Complete(A) IO[bool]
}

func MakeDeferred[A any]() IO[Deferred[A]]
```
Deferred is an interface representing a value which may not be available yet.

A deferred value can be retrieved using the `Get` function.
This value can only be set once using `Complete`.

The `Get` function blocks semantically until the Deferred is completed.

#### Example
```go
// example.go
package main

import (
  "fmt"
  c "github.com/valentinHenry/giog/io/console"
  "github.com/valentinHenry/giog/io/io"
  p "github.com/valentinHenry/giog/utils/pipes"
  v "github.com/valentinHenry/giog/utils/void"
  "time"
)

func main() {
  makeAsyncWaiter := func(waiterNb int) func(io.Deferred[int]) io.VIO {
    return func(value io.Deferred[int]) io.VIO {
      return io.Go_(
        p.Pipe2(
          c.Printf("[%d] Waiting for value\n", waiterNb).Void(),
          io.FlatMapK(func(_ v.Void) io.IO[int] { return value.Get() }),
          io.FlatMapK(func(value int) io.IO[int] { return c.Printf("[%d] Got the value: %d!\n", waiterNb, value) }),
        ),
      )
    }
  }

  give := func(toGive int) func(io.Deferred[int]) io.VIO {
    return func(value io.Deferred[int]) io.VIO {
      return p.Pipe2(
        c.Printf("Giving the value %d\n", toGive),
        io.FlatMapK(func(_ int) io.IO[bool] { return value.Complete(toGive) }),
        io.IfIOK(c.Printf("Value %d given\n", toGive), c.Printf("Failed to give %d\n", toGive)),
      ).Void()
    }
  }

  waitSomeTime := io.Delay(func() v.Void { time.Sleep(50 * time.Millisecond); return v.Void{} })

  mainLogic := p.Pipe8(
    io.MakeDeferred[int](),
    io.FlatTapK(makeAsyncWaiter(1)),
    io.FlatTapK(makeAsyncWaiter(2)),
    io.FlatTapK(AndThenTapK[Deferred[int], v.Void](waitSomeTime)),
    io.FlatTapK(give(42)),
    io.FlatTapK(makeAsyncWaiter(3)),
    io.FlatTapK(AndThenTapK[Deferred[int], v.Void](waitSomeTime)),
    io.FlatTapK(give(-1)),
    io.FlatTapK(makeAsyncWaiter(4)),
  )

  _, err := io.EvalSync(mainLogic)
  if err != nil {
    fmt.Println(err.Error())
  }
}
```
```shell
$> go run example.go
[1] Waiting for value
[2] Waiting for value
Giving the value 42
Value 42 given
[2] Got the value: 42!
[1] Got the value: 42!
[3] Waiting for value
[3] Got the value: 42!
Giving the value -1
Failed to give -1
[4] Waiting for value
[4] Got the value: 42!
```

---
### Queue
```go
type Queue[A any] interface {
  Enqueue(a A) IO[v.Void]
  TryEnqueue(a A) IO[bool]
  
  Dequeue() IO[A]
  TryDequeue() IO[o.Option[A]]
}

func BoundedQueue[A any](nb uint) IO[Queue[A]]
func UnboundedQueue[A any]() IO[Queue[A]]
func SyncQueue[A any](nb uint) IO[Queue[A]]
```

Queue is the interface of a concurrent queue

Three implementations are available:
- **BoundedQueue**: a queue which, when full, will block semantically on Enqueue(A) and return false on TryEnqueue(A)
- **UnboundedQueue**: a limitless queue
- **SyncQueue**: a queue which blocks until there is at least one reader and one writer waiting.

#### Example
// TODO

### Ref
```go
type Ref[A any] interface {
  Get() IO[A]
  Set(A) VIO
  Update(func(A) A) VIO
  TryUpdate(func(A) A) IO[bool]
  GetAndSet(A) IO[A]
  GetAndUpdate(func(A) A) IO[A]
  UpdateAndGet(func(A) A) IO[A]
}
func ModifyRef[A, B any](r Ref[A], modify func(A) (A, B)) IO[B]

func MakeRef[A any](v A) IO[Ref[A]]
```
Ref is an interface representing the reference to a value.

A Ref always references a value.

### Semaphore
```go
type Semaphore interface {
  Release() VIO
  ReleaseN(n uint) VIO
  
  Acquire() VIO
  AcquireN(n uint) VIO
  
  TryAcquire() IO[bool]
  TryAcquireN(n uint) IO[bool]
  
  Use() RIO[v.Void]
  UseN(n uint) RIO[v.Void]
}

func MakeSemaphore(size uint) IO[Semaphore]
```

 Semaphore is an interface to a non-negative amount of permits.

## Monads
### Either
### Option

## Helpers
### Functions
### Pipes
### Tuples