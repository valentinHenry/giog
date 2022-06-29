package io

import "fmt"

func _TODO[Out any]() Out {
	trace := getTrace(1)
	panic(fmt.Errorf("TODO at %s:%d\nCalled by: %s", trace.file, trace.line, trace.function))
}
