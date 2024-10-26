package data_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/pkg/data"
)

func TestMakeChecker(t *testing.T) {
	as := assert.New(t)
	fn1 := data.MakeChecker()
	as.NotNil(fn1)
	as.Nil(fn1(-1))
	as.Nil(fn1(1000))

	fn2 := data.MakeChecker(1)
	as.Nil(fn2(1))
	as.EqualError(fn2(2), fmt.Sprintf(data.ErrFixedArity, 1, 2))

	fn3 := data.MakeChecker(2, data.OrMore)
	as.Nil(fn3(5))
	as.EqualError(fn3(1), fmt.Sprintf(data.ErrMinimumArity, 2, 1))

	fn4 := data.MakeChecker(2, 7)
	as.Nil(fn4(4))
	as.EqualError(fn4(8), fmt.Sprintf(data.ErrRangedArity, 2, 7, 8))

	defer as.ExpectPanic(errors.New(data.ErrTooManyArguments))
	data.MakeChecker(1, 2, 3)
}
