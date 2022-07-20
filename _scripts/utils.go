package _scripts

import (
	"fmt"
	iofmt "github.com/valentinHenry/giog/io/fmt"
	"github.com/valentinHenry/giog/io/io"
	"github.com/valentinHenry/giog/utils/monads/either"
	v "github.com/valentinHenry/giog/utils/void"
	fio "io"
)

func typeParams(i int, prefix string) string {
	if i == 1 {
		return fmt.Sprintf("%s1", prefix)
	} else {
		return fmt.Sprintf("%s, %s%d", typeParams(i-1, prefix), prefix, i)
	}
}

func NTimes[A any](from int, until int, do func(int) io.IO[A]) io.VIO {
	return io.TailRec(from, func(curr int) io.IO[either.Either[int, v.Void]] {
		return io.If(
			curr > until,
			io.Pure(either.ToRight[int, v.Void](v.Void{})),
			io.As(do(curr), either.ToLeft[int, v.Void](curr+1)),
		)
	})
}

func PrintAllN(file fio.Writer, from int, until int, fn func(i int) string) io.VIO {
	return NTimes(from, until, func(curr int) io.IO[int] { return iofmt.Fprintln(file, fn(curr)) })
}
