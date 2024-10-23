package eval_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core/bootstrap"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/eval"
	"github.com/kode4food/ale/pkg/read"
)

func TestBasicEval(t *testing.T) {
	as := assert.New(t)

	e := bootstrap.DevNullEnvironment()
	ns := e.GetAnonymous()

	v1, err := eval.String(ns, "(if true 1 0)")
	as.Nil(err)
	as.Number(1, v1)

	v2, err := eval.String(ns, "((lambda (x) (* x 2)) 50)")
	as.Nil(err)
	as.Number(100, v2)

	v3, err := eval.String(ns, "(first (concat [1 2 3] [4 5 6]))")
	as.Nil(err)
	as.Number(1, v3)

	res, err := eval.String(ns, "(define x 99)")
	as.Number(99, res)
	as.Nil(err)
	as.Number(99, as.IsBound(ns, "x"))

	v4, err := eval.String(ns, "(and true true)")
	as.Nil(err)
	as.True(v4)
}

func TestBuiltIns(t *testing.T) {
	as := assert.New(t)

	e := bootstrap.DevNullEnvironment()
	b := e.GetAnonymous()
	ns := e.GetRoot()

	as.Nil(ns.Declare("hello").Bind(
		data.MakeProcedure(func(...data.Value) data.Value {
			return S("there")
		}, 0),
	))

	tr := read.FromString(`(hello)`)
	res, err := eval.Block(b, tr)
	as.Nil(err)
	as.String("there", res)
}
