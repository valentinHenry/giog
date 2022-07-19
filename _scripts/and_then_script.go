package _scripts

import (
	"fmt"
	"github.com/valentinHenry/giog/io/io"
	fio "io"
)

func andThenParams(i int) string {
	if i == 1 {
		return fmt.Sprint("v1 ", "IO[T1]")
	}
	return fmt.Sprintf(
		"%s, v%d IO[T%d]",
		andThenParams(i-1),
		i,
		i,
	)
}

func andThenFlatMap(curr int, until int) string {
	if curr == until {
		return fmt.Sprint("v", until)
	}

	return fmt.Sprintf(
		"_FlatMap(_trace, v%d, func (_ T%d) IO[T%d] { return %s })",
		curr,
		curr,
		until,
		andThenFlatMap(curr+1, until),
	)
}

func andThenN(curr int) string {
	return fmt.Sprintf(
		"// AndThen%d executes sequentially the %d IOs and returns the value of the last one.\nfunc AndThen%d[%s any](%s) IO[T%d] {\n\treturn _AndThen%d(getTrace(1), %s)\n}",
		curr,
		curr,
		curr,
		typeParams(curr, "T"),
		andThenParams(curr),
		curr,
		curr,
		typeParams(curr, "v"),
	)
}

func _andThenN(curr int) string {
	return fmt.Sprintf(
		"func _AndThen%d[%s any](_trace *Trace, %s) IO[T%d] {return %s}",
		curr,
		typeParams(curr, "T"),
		andThenParams(curr),
		curr,
		andThenFlatMap(1, curr),
	)
}

func writeAndThenFunctions(file fio.Writer, nbFuncs int) io.VIO {
	return PrintAllN(file, 2, nbFuncs, func(i int) string { return andThenN(i) })
}

func writeInternalAndThenFunctions(file fio.Writer, nbFuncs int) io.VIO {
	return PrintAllN(file, 2, nbFuncs, func(i int) string { return _andThenN(i) })
}
