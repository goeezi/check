package check_test

import (
	"testing"

	"github.com/goeezi/check"
	"github.com/stretchr/testify/assert"
)

func TestCatch(t *testing.T) {
	t.Parallel()

	assert.NoError(t, check.Catch(func() {}))

	assert.EqualError(t, check.Catch(func() {
		check.Must(errOops)
	}), errOops.Error())
}

func TestCatch1(t *testing.T) {
	t.Parallel()

	i, err := check.Catch1(func() int { return 42 })
	if assert.NoError(t, err) {
		assert.Equal(t, 42, i)
	}

	i, err = check.Catch1(func() int {
		check.Must(errOops)
		return 42
	})
	assert.EqualError(t, err, errOops.Error(), i)
}

func TestCatch2(t *testing.T) {
	t.Parallel()

	a, b, err := check.Catch2(func() (a, b int) { return 42, 56 })
	if assert.NoError(t, err) {
		assert.Equal(t, 42, a)
		assert.Equal(t, 56, b)
	}

	a, b, err = check.Catch2(func() (a, b int) {
		check.Must(errOops)
		return 42, 56
	})
	assert.EqualError(t, err, errOops.Error(), "%v %v", a, b)
}

func TestCatch3(t *testing.T) {
	t.Parallel()

	a, b, c, err := check.Catch3(func() (a, b, c int) { return 1, 2, 3 })
	if assert.NoError(t, err) {
		assert.Equal(t, 1, a)
		assert.Equal(t, 2, b)
		assert.Equal(t, 3, c)
	}

	a, b, c, err = check.Catch3(func() (a, b, c int) {
		check.Must(errOops)
		return 1, 2, 3
	})
	assert.EqualError(t, err, errOops.Error(), "%v %v %v", a, b, c)
}

func TestCatch4(t *testing.T) {
	t.Parallel()

	a, b, c, d, err := check.Catch4(func() (a, b, c, d int) { return 1, 2, 3, 4 })
	if assert.NoError(t, err) {
		assert.Equal(t, 1, a)
		assert.Equal(t, 2, b)
		assert.Equal(t, 3, c)
		assert.Equal(t, 4, d)
	}

	a, b, c, d, err = check.Catch4(func() (a, b, c, d int) {
		check.Must(errOops)
		return 1, 2, 3, 4
	})
	assert.EqualError(t, err, errOops.Error(), "%v %v %v %v", a, b, c, d)
}
