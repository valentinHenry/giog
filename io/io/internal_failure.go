package io

import (
	"errors"
	"fmt"
	l "github.com/valentinHenry/giog/collections/list"
	"reflect"
	"strings"
)

const (
	maxStackTrace = 50
)

type StackTrace struct {
	size   int
	traces l.List[*Trace]
}

func (st *StackTrace) AppendIfNecessary(trace *Trace) {
	if st.size > maxStackTrace {
		return
	}

	switch traces := st.traces.(type) {
	case *l.Nil[*Trace]:
		st.size = 1
		st.traces = traces.Append(trace)
		return

	case *l.Cons[*Trace]:
		if traces.Head != trace {
			st.size++
			st.traces = traces.Append(trace)
		}
		return

	default:
		panic("Impossible")
	}
}

type Cause interface {
	Cause() error
	getTrace() *StackTrace
	appendTraceIfNecessary(trace *Trace) Cause
	sPrettyPrint() string
}

type _RaisedError struct {
	cause      *error
	stackTrace *StackTrace
}

func (r *_RaisedError) Cause() error { return *r.cause }
func (r *_RaisedError) getTrace() *StackTrace {
	return r.stackTrace
}
func (r *_RaisedError) appendTraceIfNecessary(trace *Trace) Cause {
	r.stackTrace.AppendIfNecessary(trace)
	return r
}
func (r *_RaisedError) sPrettyPrint() string {
	return causeAsStr(r)
}

type _MappedCause struct {
	original   Cause
	cause      *error
	stackTrace *StackTrace
}

func (r *_MappedCause) Cause() error { return *r.cause }
func (r *_MappedCause) getTrace() *StackTrace {
	return r.stackTrace
}
func (r *_MappedCause) appendTraceIfNecessary(trace *Trace) Cause {
	r.stackTrace.AppendIfNecessary(trace)
	return r
}

func (r *_MappedCause) sPrettyPrint() string {
	return fmt.Sprintf("%s\n Mapped From: %s", causeAsStr(r), r.original.sPrettyPrint())
}

type Cancellation struct{}

func (c *Cancellation) Cause() error                          { return errors.New("cancelled") }
func (c *Cancellation) getTrace() *StackTrace                 { return &StackTrace{} }
func (c *Cancellation) appendTraceIfNecessary(_ *Trace) Cause { return c }
func (c *Cancellation) sPrettyPrint() string                  { return causeAsStr(c) }

func makeCancellationCause() Cause {
	return &Cancellation{}
}

func makeRaisedCause(_trace *Trace, cause *error) Cause {
	return &_RaisedError{cause, newStackTraceWith(_trace)}
}

func makeMappedCause(_trace *Trace, original Cause, newCause *error) Cause {
	return &_MappedCause{
		original:   original,
		cause:      newCause,
		stackTrace: newStackTraceWith(_trace),
	}
}

func newStackTraceWith(_trace *Trace) *StackTrace {
	return &StackTrace{traces: &l.Cons[*Trace]{_trace, &l.Nil[*Trace]{}}, size: 1}
}

func causeAsStr(c Cause) string {
	headStr := c.Cause().Error()

	tracesStr := ""

	traces, ok := c.getTrace().traces.(*l.Cons[*Trace])
	if !ok || traces == nil {
		return headStr
	}

	head := traces

	for ; ok; head, ok = head.Tail.(*l.Cons[*Trace]) {
		trace := head.Head

		splits := strings.Split(trace.tracedFrom, "/")
		name := splits[len(splits)-1]
		trailingDots := 40 - len(name) - 2
		if trailingDots < 0 {
			trailingDots = 0
		}
		trailingDotsStr := strings.Repeat("–", trailingDots)

		tracedFrom := fmt.Sprintf("[%s]%s", name, trailingDotsStr)

		tracesStr = fmt.Sprintf("  at %s%s(%s:%d)\n%s", tracedFrom, trace.function, trace.file, trace.line, tracesStr)
	}

	return fmt.Sprintf("%s: %s\n%s", reflect.TypeOf(c.Cause()), c.Cause().Error(), tracesStr)
}
