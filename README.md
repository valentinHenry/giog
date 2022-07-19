# GioG - IO Monad for Go

GioG is a library implementing an IO Monad for Go. It is shipped with other useful monads like `Either` and `Option`.

## IO
### Why using IO ?
GioG allow you to write code in a more functional way, without having to be a scribe monk rewriting again and again the following code:
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
### Composition Functions
### Asynchronicity

## RIO
### Basics
### Composition Functions

## Monads
### Either
### Option

## Helpers
### Functions
### Pipes
### Tuples