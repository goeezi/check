package check_test

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goeezi/check"
)

var errOops = errors.New("oops")

func TestError(t *testing.T) {
	t.Parallel()

	assert.PanicsWithError(t, errOops.Error(), func() {
		check.Must(errOops)
	})
}

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

func TestFail(t *testing.T) {
	t.Parallel()

	assert.PanicsWithValue(t, check.ErrNilError, func() {
		check.Fail(nil)
	})

	assert.PanicsWithError(t, errOops.Error(), func() {
		check.Fail(errOops)
	})
}

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

func TestHandleTransform(t *testing.T) {
	t.Parallel()

	suppressError := func(err error) error {
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

func TestMust(t *testing.T) {
	t.Parallel()

	assert.NoError(t, func() (err error) {
		defer check.Handle(&err)
		check.Must(nil)
		return
	}())
	assert.EqualError(t, func() (err error) {
		defer check.Handle(&err)
		check.Must(fmt.Errorf("oops"))
		return
	}(), "oops")
}

func TestMust1(t *testing.T) {
	t.Parallel()

	a, err := func() (a int, err error) {
		defer check.Handle(&err)
		a = check.Must1(strconv.Atoi("42"))
		return
	}()
	if assert.NoError(t, err) {
		assert.Equal(t, 42, a)
	}

	a, err = func() (a int, err error) {
		defer check.Handle(&err)
		a = check.Must1(strconv.Atoi("forty-two"))
		return
	}()
	assert.EqualError(t, err, "strconv.Atoi: parsing \"forty-two\": invalid syntax", a)
}

func divmod(a, b float64) (quo, rem float64, err error) {
	rem = math.Mod(a, b)
	quo = (a - rem) / b
	if math.IsNaN(quo) || math.IsNaN(rem) {
		err = fmt.Errorf("cannot divmod(%v, %v)", a, b)
	}
	return
}

func TestMust2(t *testing.T) {
	t.Parallel()

	a, b, err := func() (a, b float64, err error) {
		defer check.Handle(&err)
		a, b = check.Must2(divmod(42, 56))
		return
	}()
	if assert.NoError(t, err) {
		assert.EqualValues(t, 0, a)
		assert.EqualValues(t, 42, b)
	}

	a, b, err = func() (a, b float64, err error) {
		defer check.Handle(&err)
		a, b = check.Must2(divmod(0, 0))
		return
	}()
	assert.EqualError(t, err, "cannot divmod(0, 0)", "%v %v", a, b)
}

func muldivmod(a, b float64) (mul, quo, rem float64, err error) {
	mul = a * b
	rem = math.Mod(a, b)
	quo = (a - rem) / b
	if math.IsNaN(quo) || math.IsNaN(rem) {
		err = fmt.Errorf("cannot muldivmod(%v, %v)", a, b)
	}
	return
}

func TestMust3(t *testing.T) {
	t.Parallel()

	a, b, c, err := func() (a, b, c float64, err error) {
		defer check.Handle(&err)
		a, b, c = check.Must3(muldivmod(42, 56))
		return
	}()
	if assert.NoError(t, err) {
		assert.EqualValues(t, 2352, a)
		assert.EqualValues(t, 0, b)
		assert.EqualValues(t, 42, c)
	}

	a, b, c, err = func() (a, b, c float64, err error) {
		defer check.Handle(&err)
		a, b, c = check.Must3(muldivmod(0, 0))
		return
	}()
	assert.EqualError(t, err, "cannot muldivmod(0, 0)", "%v %v %v", a, b, c)
}

func analyzeTrades(prices ...float64) (open, high, low, close float64, err error) {
	if len(prices) == 0 {
		return 0, 0, 0, 0, errors.New("cannot analyze empty input")
	}
	open = prices[0]
	high = open
	low = open
	close = prices[len(prices)-1]
	for _, price := range prices {
		switch {
		case high < price:
			high = price
		case low > price:
			low = price
		}
	}
	return
}

func TestMust4(t *testing.T) {
	t.Parallel()

	o, h, l, c, err := func() (o, h, l, c float64, err error) {
		defer check.Handle(&err)
		o, h, l, c = check.Must4(analyzeTrades(3, 1, 4))
		return
	}()
	if assert.NoError(t, err) {
		assert.Equal(t, []float64{3, 4, 1, 4}, []float64{o, h, l, c})
	}

	o, h, l, c, err = func() (o, h, l, c float64, err error) {
		defer check.Handle(&err)
		o, h, l, c = check.Must4(analyzeTrades())
		return
	}()
	assert.EqualError(t, err, "cannot analyze empty input", "%v %v %v %v", o, h, l, c)
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
