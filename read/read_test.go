package read_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/read"
)

func TestFromString(t *testing.T) {
	as := assert.New(t)
	tr := read.FromString("99")
	as.NotNil(tr)
	as.Equal(I(99), tr.Car())
}
