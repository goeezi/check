package check_test

import (
	"errors"
	"fmt"
	"testing"

	goerrors "github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/goeezi/check"
)

func TestHandle(t *testing.T) {
	t.Parallel()

	assert.EqualError(t, func() (err error) {
		defer check.Handle(&err)
		check.Fail(errOops)
		return
	}(), "oops")

	assert.NoError(t, func() (err error) {
		defer check.Handle(&err)
		return
	}(), "oops")
}

func TestWrap(t *testing.T) {
	t.Parallel()

	err := func() (e error) {
		defer check.Wrap(&e, 1)
		crash := func() {
			check.Fail(errOops)
		}
		crash()
		return
	}()
	assert.EqualError(t, err, "oops")
	var werr *goerrors.Error
	require.True(t, errors.As(err, &werr))
	stk := werr.ErrorStack()
	frame := werr.StackFrames()[0]
	line, err := frame.SourceLine()
	assert.NoError(t, err)
	assert.Equal(t, "crash()", line, stk)

	assert.NoError(t, func() (e error) {
		defer check.Handle(&e)
		return
	}(), "oops")
}

func TestHandleTransform(t *testing.T) {
	t.Parallel()

	suppressError := func(error) error {
		return nil
	}

	assert.PanicsWithError(t, "oops", func() {
		defer check.Handle(nil)
		check.Fail(errOops)
	})

	assert.NotPanics(t, func() {
		defer check.Handle(nil, suppressError)
		check.Fail(errOops)
	})

	assert.NoError(t, func() (err error) {
		defer check.Handle(&err, suppressError)
		return
	}(), "oops")

	err := func() (err error) {
		defer check.Handle(&err,
			func(err error) error { return fmt.Errorf("🤨: %w", err) },
			func(err error) error { return fmt.Errorf("💥 %w 💥", err) },
		)
		check.Fail(errOops)
		return
	}()
	assert.Error(t, err, "💥 🤨: oops 💥")
	assert.Error(t, errors.Unwrap(err), "🤨: oops")
	assert.True(t, errors.Is(err, errOops))
}
