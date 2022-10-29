package check

import "errors"

var ErrNilError = errors.New("called Fail(nil)")

// Error wraps an error. The Must… family of functions use Error to wrap errors
// in calls to panic, while the Catch… family detect errors wrapped thus.
type Error struct {
	err error
}

// Error returns a string representation of e, thus implementing the error
// interface.
func (e Error) Error() string {
	return e.err.Error()
}

// Unwrap returns the wrapped error as required by the errors packages.
func (e Error) Unwrap() error {
	return e.err
}
