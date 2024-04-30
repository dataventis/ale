package encoder_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestConstants(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	i1 := e.AddConstant(S("hello"))
	i2 := e.AddConstant(I(42))
	i3 := e.AddConstant(S("hello"))
	as.Equal(i1, i3)

	c := e.Encode().Constants
	as.Equal(2, len(c))
	as.Equal(S("hello"), c[i1])
	as.Equal(I(42), c[i2])
}
