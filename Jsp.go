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
