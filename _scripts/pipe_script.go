package _scripts

import (
	"fmt"
	"github.com/valentinHenry/giog/io/io"
	fio "io"
)

func pipeFnParams(i int) string {
	if i == 1 {
		return fmt.Sprintf("fn1 func(T1) T2")
	}
	return fmt.Sprintf("%s, fn%d func (T%d) T%d", pipeFnParams(i-1), i, i, i+1)
}

func pipeFnRes(i int) string {
	if i == 1 {
		return "fn1(v)"
	}
	return fmt.Sprint("fn", i, "(", pipeFnRes(i-1), ")")
}

func pipeN(curr int) string {
	return fmt.Sprintf(
		"func Pipe%d[%s any](v T1, %s) T%d {\n\treturn %s\n}",
		curr,
		typeParams(curr+1, "T"),
		pipeFnParams(curr),
		curr+1,
		pipeFnRes(curr),
	)
}

func pipeNK(curr int) string {
	return fmt.Sprintf(
		"func Pipe%dK[%s any](%s) func (T1) T%d {\n\treturn func (v T1) T%d {\n\t\treturn %s\n\t}\n}",
		curr,
		typeParams(curr+1, "T"),
		pipeFnParams(curr),
		curr+1,
		curr+1,
		pipeFnRes(curr),
	)
}

func writePipesFunctions(file fio.Writer, nbFuncs int) io.VIO {
	printAllPipes := PrintAllN(file, 1, nbFuncs, func(i int) string { return pipeN(i) })
	printAllPipesK := PrintAllN(file, 1, nbFuncs, func(i int) string { return pipeNK(i) })

	return io.AndThen2(printAllPipes, printAllPipesK)
}
