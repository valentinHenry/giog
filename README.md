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
  p "github.com/valentinHenry/giog/pipes"
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
### Asynchronicity
### Evaluates the IO

## DataTypes

### RIO
#### Basics
#### Composition Functions
### Queue
### Ref

## Monads
### Either
### Option

## Helpers
### Functions
### Pipes
### Tuples