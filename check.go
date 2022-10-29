package check

// Fail calls Panic(Error{err}) if err is not nil, otherwise it panics.
func Fail(err error) {
	if err == nil {
		panic(ErrNilError)
	}
	panic(Error{err})
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
