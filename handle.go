package check

import (
	"math"

	"github.com/go-errors/errors"
)

// Handle, when deferred, recovers Error{err}. If any transforms are specified,
// err is transformed via err = transforms[i](err) for each transform in turn.
// Finally, Handle assigns err to *e unless e is nil, in which case it panics
// with Error{err}.
//
//	func getTotalWeight(weight, qty string) (_ float64, e error) {
//		defer Handle(&e, func(e error) error {
//			return fmt.Errorf("computing total weight: %w", e)
//		})
//		return Must1(strconv.ParseFloat(weight, 64)) *
//			float64(Must1(strconv.Atoi(qty))), nil
//	}
func Handle(e *error, transforms ...func(e error) error) {
	handle(recover(), math.MinInt, e, transforms...)
}

// Wrap behaves like Handle, but additionally wraps any returned error in
// "github.com/go-errors/errors".Error, which provides access to the stack
// trace. Use skip to drop uninteresting stack frames.
func Wrap(e *error, skip int, transforms ...func(e error) error) {
	handle(recover(), skip, e, transforms...)
}

func handle(r any, skip int, e *error, transforms ...func(e error) error) {
	if r != nil {
		if wrapped, is := r.(Error); is {
			err := wrapped.Unwrap()
			for _, transform := range transforms {
				if err = transform(err); err == nil {
					return
				}
			}
			if e == nil {
				panic(Error{err})
			}
			if skip != math.MinInt {
				err = errors.Wrap(err, 4+skip)
			}
			*e = err
			return
		}
		panic(r)
	}
}
