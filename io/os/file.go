package os

import (
	"github.com/valentinHenry/giog/io/io"
	stdos "os"
)

// TODO add a functional interface to File

func Open(name string) io.RIO[*File] {
	return OpenFile(name, O_RDONLY, 0)
}

func Create(name string) io.RIO[*File] {
	return OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
}

func OpenFile(name string, flag int, perm FileMode) io.RIO[*File] {
	return io.MakeRIO(
		io.Lift(func() (*stdos.File, error) { return stdos.OpenFile(name, flag, perm) }),
		func(file *stdos.File) io.VIO { return io.LiftV(file.Close) },
	)
}

type File = stdos.File
type FileMode = stdos.FileMode

const (
	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	O_RDONLY int = stdos.O_RDONLY // open the file read-only.
	O_WRONLY int = stdos.O_WRONLY // open the file write-only.
	O_RDWR   int = stdos.O_RDWR   // open the file read-write.
	// The remaining values may be or'ed in to control behavior.
	O_APPEND int = stdos.O_APPEND // append data to the file when writing.
	O_CREATE int = stdos.O_CREATE // create a new file if none exists.
	O_EXCL   int = stdos.O_EXCL   // used with O_CREATE, file must not exist.
	O_SYNC   int = stdos.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = stdos.O_TRUNC  // truncate regular writable file when opened.
)
