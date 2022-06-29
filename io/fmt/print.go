package fmt

import (
	"fmt"
	"github.com/valentinHenry/giog/io/io"
	fio "io"
)

func Fprintf(w fio.Writer, format string, a ...any) io.IO[int] {
	return io.Lift(func() (int, error) { return fmt.Fprintf(w, format, a...) })
}

func FprintfK(w fio.Writer, format string, a ...any) io.IO[int] {
	return io.Lift(func() (int, error) { return fmt.Fprintf(w, format, a...) })
}

func Printf(format string, a ...any) io.IO[int] {
	return io.Lift(func() (int, error) { return fmt.Printf(format, a...) })
}

func Sprintf(format string, a ...any) io.IO[string] {
	return io.Pure(fmt.Sprintf(format, a...))
}

func Fprint(w fio.Writer, a ...any) io.IO[int] {
	return io.Lift(func() (int, error) { return fmt.Fprint(w, a...) })
}

func Print(a ...any) io.IO[int] {
	return io.Lift(func() (int, error) { return fmt.Print(a) })
}

func Sprint(a ...any) io.IO[string] {
	return io.Pure(fmt.Sprint(a...))
}

func Fprintln(w fio.Writer, a ...any) io.IO[int] {
	return io.Lift(func() (int, error) { return fmt.Fprintln(w, a...) })
}

func Println(a ...any) io.IO[int] {
	return io.Lift(func() (int, error) { return fmt.Println(a...) })
}

func Sprintln(a ...any) io.IO[string] {
	return io.Pure(fmt.Sprintln(a))
}
