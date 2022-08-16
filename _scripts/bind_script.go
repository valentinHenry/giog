package _scripts

import (
	"fmt"
	"github.com/valentinHenry/giog/io/io"
	fio "io"
)

func bindParams(i int, M string) string {
	if i == 1 {
		return fmt.Sprint("v1 ", M, "[T1]")
	}
	return fmt.Sprintf(
		"%s, v%d func (%s) %s[T%d]",
		bindParams(i-1, M),
		i,
		typeParams(i-1, "T"),
		M,
		i,
	)
}

func bindIOBody(curr int, until int) string {
	if curr == until {
		return fmt.Sprintf(
			"_Map(_trace, v%d(%s), func (v%d T%d) tuples.T%d[%s] {return tuples.Of%d(%s)})",
			curr,
			typeParams(curr-1, "v"),
			curr,
			curr,
			curr,
			typeParams(curr, "T"),
			curr,
			typeParams(curr, "v"),
		)
	}

	if curr == 1 {
		return fmt.Sprintf(
			"_FlatMap(_trace, v%d, func (v%d T%d) IO[tuples.T%d[%s]] { return %s })",
			curr,
			curr,
			curr,
			until,
			typeParams(until, "T"),
			bindIOBody(curr+1, until),
		)
	}

	return fmt.Sprintf(
		"_FlatMap(_trace, v%d(%s), func (v%d T%d) IO[tuples.T%d[%s]] { return %s })",
		curr,
		typeParams(curr-1, "v"),
		curr,
		curr,
		until,
		typeParams(until, "T"),
		bindIOBody(curr+1, until),
	)
}

func bindION(curr int) string {
	return fmt.Sprintf(
		"func Bind%d[%s any](%s) IO[tuples.T%d[%s]] {\n\treturn _Bind%d(getTrace(1), %s)\n}",
		curr,
		typeParams(curr, "T"),
		bindParams(curr, "IO"),
		curr,
		typeParams(curr, "T"),
		curr,
		typeParams(curr, "v"),
	)
}

func _bindION(curr int) string {
	return fmt.Sprintf(
		"func _Bind%d[%s any](_trace *Trace, %s) IO[tuples.T%d[%s]] {return %s}",
		curr,
		typeParams(curr, "T"),
		bindParams(curr, "IO"),
		curr,
		typeParams(curr, "T"),
		bindIOBody(1, curr),
	)
}

func bindRIOBody(curr int, until int) string {
	if curr == until {
		return fmt.Sprintf(
			"_MapRIO(_trace, v%d(%s), func (v%d T%d) tuples.T%d[%s] {return tuples.Of%d(%s)})",
			curr,
			typeParams(curr-1, "v"),
			curr,
			curr,
			curr,
			typeParams(curr, "T"),
			curr,
			typeParams(curr, "v"),
		)
	}

	if curr == 1 {
		return fmt.Sprintf(
			"_FlatMapRIO(_trace, v%d, func (v%d T%d) RIO[tuples.T%d[%s]] { return %s })",
			curr,
			curr,
			curr,
			until,
			typeParams(until, "T"),
			bindRIOBody(curr+1, until),
		)
	}

	return fmt.Sprintf(
		"_FlatMapRIO(_trace, v%d(%s), func (v%d T%d) RIO[tuples.T%d[%s]] { return %s })",
		curr,
		typeParams(curr-1, "v"),
		curr,
		curr,
		until,
		typeParams(until, "T"),
		bindRIOBody(curr+1, until),
	)
}

func bindRION(curr int) string {
	return fmt.Sprintf(
		"func BindRIO%d[%s any](%s) RIO[tuples.T%d[%s]] {\n\treturn _BindRIO%d(getTrace(1), %s)\n}",
		curr,
		typeParams(curr, "T"),
		bindParams(curr, "RIO"),
		curr,
		typeParams(curr, "T"),
		curr,
		typeParams(curr, "v"),
	)
}

func _bindRION(curr int) string {
	return fmt.Sprintf(
		"func _BindRIO%d[%s any](_trace *Trace, %s) RIO[tuples.T%d[%s]] {return %s}",
		curr,
		typeParams(curr, "T"),
		bindParams(curr, "RIO"),
		curr,
		typeParams(curr, "T"),
		bindRIOBody(1, curr),
	)
}

func bindIOYBody(curr int, until int) string {
	if curr == until {
		return fmt.Sprintf(
			"v%d(%s)",
			curr,
			typeParams(curr-1, "v"),
		)
	}

	if curr == 1 {
		return fmt.Sprintf(
			"_FlatMap(_trace, v%d, func (v%d T%d) IO[T%d] { return %s })",
			curr,
			curr,
			curr,
			until,
			bindIOYBody(curr+1, until),
		)
	}

	return fmt.Sprintf(
		"_FlatMap(_trace, v%d(%s), func (v%d T%d) IO[T%d] { return %s })",
		curr,
		typeParams(curr-1, "v"),
		curr,
		curr,
		until,
		bindIOYBody(curr+1, until),
	)
}

func bindIOYN(curr int) string {
	return fmt.Sprintf(
		"func BindY%d[%s any](%s) IO[T%d] {\n\treturn _BindY%d(getTrace(1), %s)\n}",
		curr,
		typeParams(curr, "T"),
		bindParams(curr, "IO"),
		curr,
		curr,
		typeParams(curr, "v"),
	)
}

func _bindIOYN(curr int) string {
	return fmt.Sprintf(
		"func _BindY%d[%s any](_trace *Trace, %s) IO[T%d] {return %s}",
		curr,
		typeParams(curr, "T"),
		bindParams(curr, "IO"),
		curr,
		bindIOYBody(1, curr),
	)
}

func bindRIOYBody(curr int, until int) string {
	if curr == until {
		return fmt.Sprintf(
			"v%d(%s)",
			curr,
			typeParams(curr-1, "v"),
		)
	}

	if curr == 1 {
		return fmt.Sprintf(
			"_FlatMapRIO(_trace, v%d, func (v%d T%d) RIO[T%d] { return %s })",
			curr,
			curr,
			curr,
			until,
			bindRIOYBody(curr+1, until),
		)
	}

	return fmt.Sprintf(
		"_FlatMapRIO(_trace, v%d(%s), func (v%d T%d) RIO[T%d] { return %s })",
		curr,
		typeParams(curr-1, "v"),
		curr,
		curr,
		until,
		bindRIOYBody(curr+1, until),
	)
}

func bindRIOYN(curr int) string {
	return fmt.Sprintf(
		"func BindRIOY%d[%s any](%s) RIO[T%d] {\n\treturn _BindRIOY%d(getTrace(1), %s)\n}",
		curr,
		typeParams(curr, "T"),
		bindParams(curr, "RIO"),
		curr,
		curr,
		typeParams(curr, "v"),
	)
}

func _bindRIOYN(curr int) string {
	return fmt.Sprintf(
		"func _BindRIOY%d[%s any](_trace *Trace, %s) RIO[T%d] {return %s}",
		curr,
		typeParams(curr, "T"),
		bindParams(curr, "RIO"),
		curr,
		bindRIOYBody(1, curr),
	)
}

func bindOptionBody(curr int, until int) string {
	if curr == until {
		return fmt.Sprintf(
			"Map(v%d(%s), func (v%d T%d) tuples.T%d[%s] {return tuples.Of%d(%s)})",
			curr,
			typeParams(curr-1, "v"),
			curr,
			curr,
			curr,
			typeParams(curr, "T"),
			curr,
			typeParams(curr, "v"),
		)
	}

	if curr == 1 {
		return fmt.Sprintf(
			"FlatMap(v%d, func (v%d T%d) Option[tuples.T%d[%s]] { return %s })",
			curr,
			curr,
			curr,
			until,
			typeParams(until, "T"),
			bindOptionBody(curr+1, until),
		)
	}

	return fmt.Sprintf(
		"FlatMap(v%d(%s), func (v%d T%d) Option[tuples.T%d[%s]] { return %s })",
		curr,
		typeParams(curr-1, "v"),
		curr,
		curr,
		until,
		typeParams(until, "T"),
		bindOptionBody(curr+1, until),
	)
}

func bindOptionN(curr int) string {
	return fmt.Sprintf(
		"func Bind%d[%s any](%s) Option[tuples.T%d[%s]] {return %s}",
		curr,
		typeParams(curr, "T"),
		bindParams(curr, "Option"),
		curr,
		typeParams(curr, "T"),
		bindOptionBody(1, curr),
	)
}

func bindOptionYBody(curr int, until int) string {
	if curr == until {
		return fmt.Sprintf(
			"v%d(%s)",
			curr,
			typeParams(curr-1, "v"),
		)
	}

	if curr == 1 {
		return fmt.Sprintf(
			"FlatMap(v%d, func (v%d T%d) Option[T%d] { return %s })",
			curr,
			curr,
			curr,
			until,
			bindOptionYBody(curr+1, until),
		)
	}

	return fmt.Sprintf(
		"FlatMap(v%d(%s), func (v%d T%d) Option[T%d] { return %s })",
		curr,
		typeParams(curr-1, "v"),
		curr,
		curr,
		until,
		bindOptionYBody(curr+1, until),
	)
}

func bindOptionYN(curr int) string {
	return fmt.Sprintf(
		"func BindY%d[%s any](%s) Option[T%d] {return %s}",
		curr,
		typeParams(curr, "T"),
		bindParams(curr, "Option"),
		curr,
		bindOptionYBody(1, curr),
	)
}

func writeIoBindFunctions(file fio.Writer, nbFuncs int) io.VIO {
	return io.AndThen4(
		PrintAllN(file, 2, nbFuncs, func(i int) string { return bindION(i) }),
		PrintAllN(file, 2, nbFuncs, func(i int) string { return bindRION(i) }),
		PrintAllN(file, 2, nbFuncs, func(i int) string { return bindIOYN(i) }),
		PrintAllN(file, 2, nbFuncs, func(i int) string { return bindRIOYN(i) }),
	)
}
func writeInternalIoBindFunctions(file fio.Writer, nbFuncs int) io.VIO {
	return io.AndThen4(
		PrintAllN(file, 2, nbFuncs, func(i int) string { return _bindION(i) }),
		PrintAllN(file, 2, nbFuncs, func(i int) string { return _bindRION(i) }),
		PrintAllN(file, 2, nbFuncs, func(i int) string { return _bindIOYN(i) }),
		PrintAllN(file, 2, nbFuncs, func(i int) string { return _bindRIOYN(i) }),
	)
}

func writeOptionBindFunctions(file fio.Writer, nbFuncs int) io.VIO {
	return io.AndThen2(
		PrintAllN(file, 2, nbFuncs, func(i int) string { return bindOptionN(i) }),
		PrintAllN(file, 2, nbFuncs, func(i int) string { return bindOptionYN(i) }),
	)
}
