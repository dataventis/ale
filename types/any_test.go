package types_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	as := assert.New(t)

	as.Equal("any", types.Any.Name())
	as.True(types.Accepts(types.Any, types.Lambda))
	as.True(types.Accepts(types.Any, types.Number))
	as.True(types.Accepts(types.Any, types.Any))
}
