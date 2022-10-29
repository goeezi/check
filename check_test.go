package check_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goeezi/check"
)

func TestFail(t *testing.T) {
	t.Parallel()

	assert.PanicsWithValue(t, check.ErrNilError, func() {
		check.Fail(nil)
	})

	assert.PanicsWithError(t, errOops.Error(), func() {
		check.Fail(errOops)
	})
}

func TestPass(t *testing.T) {
	t.Parallel()

	var log []string

	report := func(format string, args ...any) {
		log = append(log, fmt.Sprintf(format, args...))
	}

	logPanic := func() {
		if r := check.Pass(recover()); r != nil {
			report(fmt.Sprintf("error: %v", r))
			panic(r)
		}
	}

	assert.NoError(t, func() (err error) {
		defer check.Handle(&err)
		defer logPanic()
		return nil
	}())
	assert.Empty(t, log)

	assert.Error(t, func() (err error) {
		defer check.Handle(&err)
		defer logPanic()
		check.Fail(errOops)
		return nil
	}(), "oops")
	assert.Empty(t, log)

	assert.PanicsWithValue(t, 42, func() {
		var err error
		defer check.Handle(&err)
		defer logPanic()
		panic(42)
	})
	assert.Equal(t, []string{"error: 42"}, log)
}
