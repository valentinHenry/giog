package _scripts

import (
	"fmt"
	iofmt "github.com/valentinHenry/giog/io/fmt"
	"github.com/valentinHenry/giog/io/io"
	fio "io"
)

func fnN(i int) string {
	types := typeParams(i, "T")
	return fmt.Sprint("type Fn", i, "[", types, ", Out any] func (", types, ")", " Out")
}

func tupled(i int) string {
	types := typeParams(i, "T")
	values := typeParams(i, "v")

	return fmt.Sprintf(
		"func (function Fn%d[%s, Out]) Tupled(tp tuples.T%d[%s]) Out {\n\t%s := tp.Values()\n\treturn function(%s)\n}",
		i,
		types,
		i,
		types,
		values,
		values,
	)
}

func writeFunctionsFunctions(file fio.Writer, nbFuncs int) io.VIO {
	printAllFn := PrintAllN(file, 1, nbFuncs, fnN)
	printAllTupled := PrintAllN(file, 2, nbFuncs, tupled)

	return io.AndThen3(
		iofmt.Fprintln(file, "func Identity[A any](a A) A {\n\treturn a}\n"),
		printAllFn,
		printAllTupled,
	)
}
