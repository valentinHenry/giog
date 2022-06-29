package io

import (
	"runtime"
)

type Trace struct {
	file       string
	line       int
	function   string
	tracedFrom string
}

func getTrace(skip int) *Trace {
	// skip + 1 because we want to skip getTrace
	tracedFrom, currFrame := getFrames(skip + 1)

	return &Trace{
		file:       currFrame.File,
		line:       currFrame.Line,
		function:   currFrame.Function,
		tracedFrom: tracedFrom.Function,
	}
}

func getFrames(skipFrames int) (runtime.Frame, runtime.Frame) {
	programCounters := make([]uintptr, 2)
	// +2 because we skips getFrames
	n := runtime.Callers(skipFrames+1, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	frame2 := runtime.Frame{Function: "unknown"}

	frames := runtime.CallersFrames(programCounters[:n])

	frameCandidate, _ := frames.Next()
	frame = frameCandidate
	frameCandidate, _ = frames.Next()
	frame2 = frameCandidate

	return frame, frame2
}
