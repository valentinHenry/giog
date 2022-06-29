package time

import (
	"github.com/valentinHenry/giog/io/io"
	stdtime "time"
)

func Now() io.IO[Time] {
	return io.Delay(func() Time { return stdtime.Now() })
}

type Time = stdtime.Time
