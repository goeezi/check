package check_test

import (
	"errors"
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
