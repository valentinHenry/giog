package io

import (
	v "github.com/valentinHenry/giog/void"
)

type IO[A any] interface {
	run(*Universe) runResult[A]
	UnsafeRun() (A, error)
	Void() IO[v.Void]
}

type VIO = IO[v.Void]
