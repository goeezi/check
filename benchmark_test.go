package check_test

import (
	"testing"

	"github.com/goeezi/check"
)

func failer() (i int, err error) {
	return 0, errOops
}

func succeeder() (i int, err error) {
	return 42, nil
}

func call(f func() (int, error)) int {
	i, err := f()
	if err != nil {
		return -1
	}
	return i
}

func catch(f func() (int, error)) (i int) {
	if check.Catch(func() {
		i = check.Must1(f())
	}) != nil {
		i = -1
	}
	return
}

func handle(f func() (int, error)) int {
	var err error
	defer check.Handle(&err)
	return check.Must1(f())
}

func handleTransform(f func() (int, error)) (i int) {
	var err error
	defer check.Handle(&err, func(err error) error {
		i = -1
		return nil
	})
	return check.Must1(f())
}

func BenchmarkFailureConventional(b *testing.B) {
	for i := 0; i < b.N; i++ {
		call(failer)
	}
}

func BenchmarkFailureCatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		catch(failer)
	}
}

func BenchmarkFailureHandle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		handle(failer)
	}
}

func BenchmarkFailureHandleTransform(b *testing.B) {
	for i := 0; i < b.N; i++ {
		handleTransform(failer)
	}
}

func BenchmarkSuccessConventional(b *testing.B) {
	for i := 0; i < b.N; i++ {
		call(succeeder)
	}
}

func BenchmarkSuccessCatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		catch(succeeder)
	}
}

func BenchmarkSuccessHandle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		handle(succeeder)
	}
}

func BenchmarkSuccessHandleTransform(b *testing.B) {
	for i := 0; i < b.N; i++ {
		handleTransform(succeeder)
	}
}
