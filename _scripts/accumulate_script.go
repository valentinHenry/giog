package _scripts

import (
	"fmt"
	"github.com/valentinHenry/giog/io/io"
	fio "io"
)

func accumulateParams(i int, M string) string {
	if i == 1 {
		return fmt.Sprint("v1 ", M, "[T1]")
	}
	return fmt.Sprintf(
		"%s, v%d %s[T%d]",
		accumulateParams(i-1, M),
		i,
		M,
		i,
	)
}

func accumulateIOBody(curr int, until int) string {
	if curr == until {
		return fmt.Sprintf(
			"_Map(_trace, v%d, func (v%d T%d) tuples.T%d[%s] {return tuples.Of%d(%s)})",
			curr,
			curr,
			curr,
			curr,
			typeParams(curr, "T"),
			curr,
			typeParams(curr, "v"),
		)
	}

	return fmt.Sprintf(
		"_FlatMap(_trace, v%d, func (v%d T%d) IO[tuples.T%d[%s]] { return %s })",
		curr,
		curr,
		curr,
		until,
		typeParams(until, "T"),
		accumulateIOBody(curr+1, until),
	)
}

func accumulateION(curr int) string {
	return fmt.Sprintf(
		"//Accumulate%d evaluates sequentially the %d IOs and returns a tuple made of all IOs evaluated values.\nfunc Accumulate%d[%s any](%s) IO[tuples.T%d[%s]] {\n\treturn _Accumulate%d(getTrace(1), %s)\n}",
		curr,
		curr,
		curr,
		typeParams(curr, "T"),
		accumulateParams(curr, "IO"),
		curr,
		typeParams(curr, "T"),
		curr,
		typeParams(curr, "v"),
	)
}

func _accumulateION(curr int) string {
	return fmt.Sprintf(
		"func _Accumulate%d[%s any](_trace *Trace, %s) IO[tuples.T%d[%s]] {return %s}",
		curr,
		typeParams(curr, "T"),
		accumulateParams(curr, "IO"),
		curr,
		typeParams(curr, "T"),
		accumulateIOBody(1, curr),
	)
}

func accumulateRIOBody(curr int, until int) string {
	if curr == until {
		return fmt.Sprintf(
			"_MapRIO(_trace, v%d, func (v%d T%d) tuples.T%d[%s] {return tuples.Of%d(%s)})",
			curr,
			curr,
			curr,
			curr,
			typeParams(curr, "T"),
			curr,
			typeParams(curr, "v"),
		)
	}

	return fmt.Sprintf(
		"_FlatMapRIO(_trace, v%d, func (v%d T%d) RIO[tuples.T%d[%s]] { return %s })",
		curr,
		curr,
		curr,
		until,
		typeParams(until, "T"),
		accumulateRIOBody(curr+1, until),
	)
}

func accumulateRION(curr int) string {
	return fmt.Sprintf(
		"//AccumulateRIO%d evaluates sequentially the %d RIOs and returns a tuple made of all RIOs evaluated values.\nfunc AccumulateRIO%d[%s any](%s) RIO[tuples.T%d[%s]] {\n\treturn _AccumulateRIO%d(getTrace(1), %s)\n}",
		curr,
		curr,
		curr,
		typeParams(curr, "T"),
		accumulateParams(curr, "RIO"),
		curr,
		typeParams(curr, "T"),
		curr,
		typeParams(curr, "v"),
	)
}

func _accumulateRION(curr int) string {
	return fmt.Sprintf(
		"func _AccumulateRIO%d[%s any](_trace *Trace, %s) RIO[tuples.T%d[%s]] {return %s}",
		curr,
		typeParams(curr, "T"),
		accumulateParams(curr, "RIO"),
		curr,
		typeParams(curr, "T"),
		accumulateRIOBody(1, curr),
	)
}

func accumulateOptionBody(curr int, until int) string {
	if curr == until {
		return fmt.Sprintf(
			"Map(v%d, func (v%d T%d) tuples.T%d[%s] {return tuples.Of%d(%s)})",
			curr,
			curr,
			curr,
			curr,
			typeParams(curr, "T"),
			curr,
			typeParams(curr, "v"),
		)
	}

	return fmt.Sprintf(
		"FlatMap(v%d, func (v%d T%d) Option[tuples.T%d[%s]] { return %s })",
		curr,
		curr,
		curr,
		until,
		typeParams(until, "T"),
		accumulateOptionBody(curr+1, until),
	)
}

func accumulateOptionN(curr int) string {
	return fmt.Sprintf(
		"func Accumulate%d[%s any](%s) Option[tuples.T%d[%s]] {return %s}",
		curr,
		typeParams(curr, "T"),
		accumulateParams(curr, "Option"),
		curr,
		typeParams(curr, "T"),
		accumulateOptionBody(1, curr),
	)
}

func writeIoAccumulateFunctions(file fio.Writer, nbFuncs int) io.VIO {
	return io.AndThen2(
		PrintAllN(file, 2, nbFuncs, func(i int) string { return accumulateION(i) }),
		PrintAllN(file, 2, nbFuncs, func(i int) string { return accumulateRION(i) }),
	)
}

func writeInternalIoAccumulateFunctions(file fio.Writer, nbFuncs int) io.VIO {
	return io.AndThen2(
		PrintAllN(file, 2, nbFuncs, func(i int) string { return _accumulateION(i) }),
		PrintAllN(file, 2, nbFuncs, func(i int) string { return _accumulateRION(i) }),
	)
}

func writeOptionAccumulateFunctions(file fio.Writer, nbFuncs int) io.VIO {
	return PrintAllN(file, 2, nbFuncs, func(i int) string { return accumulateOptionN(i) })
}
