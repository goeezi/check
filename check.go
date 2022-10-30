package check

import "fmt"

// Fail panics Error{err} if err is not nil, otherwise it panics with an
// internal error that isn't treated specially by Handle, Catch[N] and Pass.
func Fail(err error) {
	if err == nil {
		panic(ErrNilError)
	}
	panic(Error{err})
}

// Fail panics Error{fmt.Errorf(format, args...)}.
func Failf(format string, args ...any) {
	panic(Error{fmt.Errorf(format, args...)})
}

// Pass returns r unless it is a check.Error, in which case it re-panics r.
// Typical usage:
//
//	defer func() {
//		if r := check.Pass(recover()); r != nil {
//			// We only get here if recover() returns something other than
//			// check.Error.
//		}
//	}()
func Pass(r any) any {
	if _, is := r.(Error); is {
		panic(r)
	}
	return r
}
