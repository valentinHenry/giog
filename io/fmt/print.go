package fmt

import (
	"fmt"
	. "github.com/valentinHenry/giog/io/io"
	"io"
)

func Fprintf(w io.Writer, format string, a ...any) IO[int] {
	return Lift(func() (int, error) { return fmt.Fprintf(w, format, a...) })
}

func FprintfK(w io.Writer, format string, a ...any) IO[int] {
	return Lift(func() (int, error) { return fmt.Fprintf(w, format, a...) })
}

func Printf(format string, a ...any) IO[int] {
	return Lift(func() (int, error) { return fmt.Printf(format, a...) })
}

func Sprintf(format string, a ...any) IO[string] {
	return Pure(fmt.Sprintf(format, a...))
}

func Fprint(w io.Writer, a ...any) IO[int] {
	return Lift(func() (int, error) { return fmt.Fprint(w, a...) })
}

func Print(a ...any) IO[int] {
	return Lift(func() (int, error) { return fmt.Print(a) })
}

func Sprint(a ...any) IO[string] {
	return Pure(fmt.Sprint(a...))
}

func Fprintln(w io.Writer, a ...any) IO[int] {
	return Lift(func() (int, error) { return fmt.Fprintln(w, a...) })
}

func Println(a ...any) IO[int] {
	return Lift(func() (int, error) { return fmt.Println(a...) })
}

func Sprintln(a ...any) IO[string] {
	return Pure(fmt.Sprintln(a))
}
