package time

import (
	. "github.com/valentinHenry/giog/io/io"
	stdtime "time"
)

func Now() IO[Time] {
	return Delay(func() Time { return stdtime.Now() })
}

type Time = stdtime.Time
