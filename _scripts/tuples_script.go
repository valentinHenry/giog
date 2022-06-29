package _scripts

import (
	"fmt"
	"github.com/valentinHenry/giog/io/io"
	fio "io"
)

func tuplesGetter(i int) string {
	if i == 1 {
		return "\tV1() P1"
	}
	return fmt.Sprintf("%s\n\tV%d() P%d", tuplesGetter(i-1), i, i)
}

func tupleInterface(curr int) string {
	params := typeParams(curr, "P")
	return fmt.Sprintf(
		"type T%d[%s any] interface {\n\tValues() (%s)\n%s\n}",
		curr,
		params,
		params,
		tuplesGetter(curr),
	)
}

func tupleImplTypeValue(i int) string {
	if i == 1 {
		return "\t_1 P1"
	}
	return fmt.Sprintf("%s\n\t_%d P%d", tupleImplTypeValue(i-1), i, i)
}

func tupleImplGetters(curr int, i int) string {
	if i == 1 {
		return fmt.Sprintf(
			"func (t tuple%dImpl[%s]) V1() P1 {\n\treturn t._1\n}",
			curr,
			typeParams(curr, "P"),
		)
	}
	return fmt.Sprintf(
		"%s\nfunc (t tuple%dImpl[%s]) V%d() P%d {\n\treturn t._%d\n}",
		tupleImplGetters(curr, i-1),
		curr,
		typeParams(curr, "P"),
		i,
		i,
		i,
	)
}

func tupleImpl(curr int) string {
	impl := fmt.Sprintf(
		"type tuple%dImpl[%s any] struct {\n%s\n}",
		curr,
		typeParams(curr, "P"),
		tupleImplTypeValue(curr),
	)

	valueFunc := fmt.Sprintf(
		"func (t tuple%dImpl[%s]) Values() (%s) {\n\treturn %s\n}",
		curr,
		typeParams(curr, "P"),
		typeParams(curr, "P"),
		typeParams(curr, "t._"),
	)

	implGetter := tupleImplGetters(curr, curr)

	return fmt.Sprint(impl, "\n", valueFunc, "\n", implGetter)
}

func tupleOfParams(i int) string {
	if i == 1 {
		return "v1 P1"
	}
	return fmt.Sprintf("%s, v%d P%d", tupleOfParams(i-1), i, i)
}

func tupleOf(curr int) string {
	return fmt.Sprintf(
		"func Of%d[%s any](%s) T%d[%s] {\n\treturn tuple%dImpl[%s]{%s}\n}",
		curr,
		typeParams(curr, "P"),
		tupleOfParams(curr),
		curr,
		typeParams(curr, "P"),
		curr,
		typeParams(curr, "P"),
		typeParams(curr, "v"),
	)
}

func writeTuplesFunctions(file fio.Writer, nbFuncs int) io.VIO {
	printAllInterfaces := PrintAllN(file, 1, nbFuncs, tupleInterface)
	printAllImpl := PrintAllN(file, 1, nbFuncs, tupleImpl)
	printAllOf := PrintAllN(file, 1, nbFuncs, tupleOf)

	return io.AndThen3(
		printAllInterfaces,
		printAllImpl,
		printAllOf,
	)
}
