# GioG - IO Monad for Go

GioG is a library implementing an IO Monad for Go. It is shipped with other useful monads like `Either` and `Option`.

## IO
### Why using IO ?
GioG allow you to write code in a more functional way by encapsulating side effects and by removing the need to be rewriting again and again the following code:
```go
result, err := runSomething()
if err != nil {
  // Do something with the error
  return nil, err
}
...
return something, nil
```

### Basics

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

### Composition Functions
#### Pipe
#### AndThen
#### Accumulate
#### Bind
#### BindY
### Asynchronicity
### Evaluates the IO

## DataTypes

### RIO
#### Basics
#### Composition Functions

---
### CountDownLatch
```go
type CountDownLatch interface {
  Release() VIO
  Await() VIO
}
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
  r "github.com/valentinHenry/refined"
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
    io.MakeCountDownLatch(r.RefineUnsafe[int, r.Positive](3)),
    io.FlatTapK(func(latch io.CountDownLatch) io.VIO { return io.Fork_(makeWaiter(latch, 1)) }),
    io.FlatTapK(func(latch io.CountDownLatch) io.VIO { return io.Fork_(makeWaiter(latch, 2)) }),
    io.FlatTapK(releaseLatch),
    io.FlatTapK(func(latch io.CountDownLatch) io.VIO { return io.Fork_(makeWaiter(latch, 3)) }),
    io.FlatTapK(releaseLatch),
    io.FlatTapK(func(latch io.CountDownLatch) io.VIO { return io.Fork_(makeWaiter(latch, 4)) }),
    io.FlatTapK(releaseLatch),
    io.FlatTapK(func(latch io.CountDownLatch) io.VIO { return io.Fork_(makeWaiter(latch, 4)) }),
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
  r "github.com/valentinHenry/refined"
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
    io.MakeCyclicBarrier(r.RefineUnsafe[int, r.Positive](3)),
    io.FlatTapK(func(b io.CyclicBarrier) io.VIO { return io.Fork_(makeWaiter(b, 1)) }),
    io.FlatTapK(func(b io.CyclicBarrier) io.VIO { return io.Fork_(makeWaiter(b, 2)) }),
    io.FlatTapK(func(b io.CyclicBarrier) io.VIO { return io.Fork_(makeWaiter(b, 3)) }),
    io.FlatTapK(func(b io.CyclicBarrier) io.VIO { return io.Fork_(makeWaiter(b, 4)) }),
    io.FlatTapK(func(b io.CyclicBarrier) io.VIO { return io.Fork_(makeWaiter(b, 5)) }),
    io.FlatTapK(func(b io.CyclicBarrier) io.VIO { return io.Fork_(makeWaiter(b, 6)) }),
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
      return io.Fork_(
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
    io.FlatTapK(func(_ io.Deferred[int]) io.VIO { return waitSomeTime }),
    io.FlatTapK(give(42)),
    io.FlatTapK(makeAsyncWaiter(3)),
    io.FlatTapK(func(_ io.Deferred[int]) io.VIO { return waitSomeTime }),
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


### Ref
### Semaphore

## Monads
### Either
### Option

## Helpers
### Functions
### Pipes
### Tuples