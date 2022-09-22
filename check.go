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

// Catch returns err if calling work panics with Error{err}, otherwise it
// returns nil.
//
//	return check.Catch(func() {
//		check.Must1(fmt.Println("Hello, World!")
//		check.Must1(fmt.Println("¡Hola, Mundo!")
//		check.Must1(fmt.Println("你好，世界!")
//		check.Must1(fmt.Println("Привет, мир!")
//	})
func Catch(work func(), transforms ...func(e error) error) (e error) {
	defer Handle(&e, transforms...)
	work()
	return
}

// Catch1 returns _, err if calling work panics with Error{err}, otherwise it
// returns t, nil.
//
//	func getTotalWeight(weight, qty string) (float64, error) {
//		return Catch1(func() float64 {
//			return Must1(strconv.ParseFloat(weight, 64)) *
//				float64(Must1(strconv.Atoi(qty)))
//		})
//	}
func Catch1[T any](work func() T, transforms ...func(e error) error) (t T, e error) {
	defer Handle(&e, transforms...)
	t = work()
	return
}

// Catch2 returns _, _, err if calling work panics with Error{err}, otherwise it
// returns t1, t2, nil. See Catch1 for a related example.
func Catch2[T1, T2 any](
	work func() (T1, T2),
	transforms ...func(e error) error,
) (t1 T1, t2 T2, e error) {
	defer Handle(&e, transforms...)
	t1, t2 = work()
	return
}

// Catch4 returns _, _, err if calling work panics with Error{err}, otherwise it
// returns t1, t2, t3, nil. See Catch1 for a related example.
func Catch3[T1, T2, T3 any](
	work func() (T1, T2, T3),
	transforms ...func(e error) error,
) (t1 T1, t2 T2, t3 T3, e error) {
	defer Handle(&e, transforms...)
	t1, t2, t3 = work()
	return
}

// Catch4 returns _, _, _, err if calling work panics with Error{err}, otherwise
// it returns t1, t2, t3, t4 nil. See Catch1 for a related example.
func Catch4[T1, T2, T3, T4 any](
	work func() (T1, T2, T3, T4),
	transforms ...func(e error) error,
) (t1 T1, t2 T2, t3 T3, t4 T4, e error) {
	defer Handle(&e, transforms...)
	t1, t2, t3, t4 = work()
	return
}

// Fail calls Panic(Error{err}) if err is not nil, otherwise it calls
// panic(NoError).
func Fail(err error) {
	if err == nil {
		panic(ErrNilError)
	}
	panic(Error{err})
}

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
	if r := recover(); r != nil {
		if wrapper, is := r.(Error); is {
			err := wrapper.Unwrap()
			for _, transform := range transforms {
				if err = transform(err); err == nil {
					return
				}
			}
			if e == nil {
				panic(Error{err})
			}
			*e = err
			return
		}
		panic(r)
	}
}

// Must calls panic(Error{err}) if err is not nil.
func Must(err error) {
	if err != nil {
		panic(Error{err})
	}
}

// Must1 returns t if err is nil, otherwise it calls panic(Error{err}).
//
//	price := check.Must1(strconv.ParseFloat(unitPrice, 64)) *
//		check.Must1(strconv.ParseFloat(qty, 64))
func Must1[T any](t T, err error) T {
	if err != nil {
		panic(Error{err})
	}
	return t
}

// Must2 returns t1, t2 if err is nil, otherwise it calls panic(Error{err}).
//
//	// MulDiv's third return value is an error if x = y = 0.
//	prod, quo := check.Must2(MulDiv(x, y))
func Must2[T1, T2 any](t1 T1, t2 T2, err error) (T1, T2) {
	Must(err)
	return t1, t2
}

// Must3 returns t1, t2, t3 if err is nil, otherwise it calls panic(Error{err}).
//
//	// MulDivRem's fourth return value is an error if x = y = 0.
//	prod, quo, rem := check.Must3(MulDivRem(x, y))
func Must3[T1, T2, T3 any](t1 T1, t2 T2, t3 T3, err error) (T1, T2, T3) {
	Must(err)
	return t1, t2, t3
}

// Must4 returns t1, t2, t3, t4 if err is nil, otherwise it calls panic(Error{err}).
//
//	// AnalyzeTrades's fifth return value is an error if prices is empty.
//	open, high, low, close := check.Must4(AnalyzeTrades(prices))
func Must4[T1, T2, T3, T4 any](t1 T1, t2 T2, t3 T3, t4 T4, err error) (T1, T2, T3, T4) {
	Must(err)
	return t1, t2, t3, t4
}

// Pass returns r unless it is a check.Error, in which case it re-panics.
// Typical usage:
//
//	defer func() {
//	    if r := check.Pass(recover()); r != nil {
//	        // If recover() returns Error, it passes through (re-panics) without
//	        // getting here.
//	    }
//	}()
func Pass(r any) any {
	if _, is := r.(Error); is {
		panic(r)
	}
	return r
}
