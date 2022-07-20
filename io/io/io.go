package io

import (
	v "github.com/valentinHenry/giog/utils/void"
)

type IO[A any] interface {
	run(*Universe) (IO[A], *A, Cause)
	UnsafeRun() (A, error)
	Void() IO[v.Void]
}

type VIO = IO[v.Void]
