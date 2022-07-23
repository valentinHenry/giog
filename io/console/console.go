package console

import (
	"bufio"
	stdfmt "fmt"
	"github.com/valentinHenry/giog/io/fmt"
	. "github.com/valentinHenry/giog/io/io"
	"os"
)

func ReadLn() IO[string] {
	return Lift(func() (string, error) {
		var str string
		_, err := stdfmt.Scanln(&str)
		return str, err
	})
}

func ReadRune() IO[rune] {
	return Lift(func() (rune, error) {
		r, _, err := bufio.NewReader(os.Stdin).ReadRune()
		return r, err
	})
}

//func Printf(format string, a ...any) IO[int] {
//	return fmt.Printf(format, a...)
//}

func Print(a ...any) IO[int] {
	return fmt.Print(a...)
}

func Println(a ...any) IO[int] {
	return fmt.Println(a...)
}
