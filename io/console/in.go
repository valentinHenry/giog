package console

import (
	"bufio"
	"fmt"
	"github.com/valentinHenry/giog/io/io"
	"os"
)

func ReadLn() io.IO[string] {
	return io.Lift(func() (string, error) {
		var str string
		_, err := fmt.Scanln(&str)
		return str, err
	})
}

func ReadRune() io.IO[rune] {
	return io.Lift(func() (rune, error) {
		r, _, err := bufio.NewReader(os.Stdin).ReadRune()
		return r, err
	})
}
