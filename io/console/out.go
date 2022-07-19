package console

import (
	"github.com/valentinHenry/giog/io/fmt"
	"github.com/valentinHenry/giog/io/io"
)

func Printf(format string, a ...any) io.IO[int] {
	return fmt.Printf(format, a...)
}

func Print(a ...any) io.IO[int] {
	return fmt.Print(a...)
}

func Println(a ...any) io.IO[int] {
	return fmt.Println(a...)
}
