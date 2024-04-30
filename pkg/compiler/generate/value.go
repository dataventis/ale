package generate

import (
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/pkg/compiler/encoder"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/macro"
)

var consSym = env.RootSymbol("cons")

// Value encodes an expression
func Value(e encoder.Encoder, v data.Value) {
	ns := e.Globals()
	switch expanded := macro.Expand(ns, v).(type) {
	case data.Sequence:
		Sequence(e, expanded)
	case data.Pair:
		Pair(e, expanded)
	case data.Symbol:
		ReferenceSymbol(e, expanded)
	case data.Keyword, data.Number, data.Bool, data.Procedure:
		Literal(e, expanded)
	default:
		panic(debug.ProgrammerError("unknown value type: %s", v))
	}
}

// Pair encodes a pair
func Pair(e encoder.Encoder, c data.Pair) {
	f := resolveBuiltIn(e, consSym)
	args := data.Vector{c.Car(), c.Cdr()}
	callStatic(e, f, args)
}

func resolveBuiltIn(e encoder.Encoder, sym data.Symbol) data.Procedure {
	ge := e.Globals().Environment()
	root := ge.GetRoot()
	res := env.MustResolveValue(root, sym)
	return res.(data.Procedure)
}
