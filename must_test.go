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
