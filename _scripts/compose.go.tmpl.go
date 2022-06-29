package _scripts

import (
	"github.com/valentinHenry/giog/io/io"
	fio "io"
)

//TODO
//func composeParams(i int) string {
//	if i == 1 {
//		return "fn1 Kleisli[T1, T2]"
//	}
//	return fmt.Sprintf("%s, fn%d Kleisli[T%d, T%d]",
//		composeParams(i-1),
//		i,
//		i,
//		i+1,
//	)
//}
//
//func composeFlatMapK(i int) string {
//	if i == 1 {
//		return "\t\tFlatMapK(fn1),"
//	}
//	return fmt.Sprintf("%s\n\t\tFlatMapK(fn%d),",
//		composeFlatMapK(i-1),
//		i,
//	)
//}
//
//func composeN(curr int, M string) string {
//	return fmt.Sprintf(
//		"func Compose%d[%s any](%s) Kleisli[%s[T1], T%d] {\n\treturn Pipe%d(\n%s\n\t)\n}",
//		curr,
//		typeParams(curr+1, "T"),
//		composeParams(curr),
//		M,
//		curr+1,
//		curr,
//		composeFlatMapK(curr),
//	)
//}
//func xcomposeN(curr int, M string) string {
//	return ""
//}
//
//func composeNK(curr int, M string) string {
//	return fmt.Sprintf(
//		"func Compose%d[%s any](%s) Kleisli[%s[T1], T%d] {\n\treturn Pipe%d(\n%s\n\t)\n}",
//		curr,
//		typeParams(curr+1, "T"),
//		composeParams(curr),
//		M,
//		curr+1,
//		curr,
//		composeFlatMapK(curr),
//	)
//}
//func xcomposeNK(curr int, M string) string {
//	return ""
//}

func WriteComposeFile(file fio.Writer, nbFuncs int, M string, pckg string) io.VIO {
	return io.Void()
}
