package asm_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/core/asm"
)

func TestConstant(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(
		`(asm* .const ("this is a list" 1 2 3))`,
		L(S("this is a list"), I(1), I(2), I(3)),
	)

	as.EvalTo(
		`(asm* .const 1 .const 2 .const 3 add add)`,
		I(6),
	)
}

func TestOutOfScopeError(t *testing.T) {
	as := assert.New(t)
	as.PanicWith(`
		(asm*
			.push-locals
			.local wont-be-found :val
			.const "hello"
			store wont-be-found
			.pop-locals
			load wont-be-found)
    `, fmt.Errorf(asm.ErrUnexpectedName, "wont-be-found"))
}

func TestLocalScopeError(t *testing.T) {
	as := assert.New(t)
	as.PanicWith(`
		(asm*
			.pop-locals
			.local hello :val)
	`, errors.New(encoder.ErrNoLocalScope))
}

func TestEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(asm* .eval (+ 1 2))`, I(3))
	as.EncodesAs(isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(1),
		isa.Const.New(0),
		isa.Call.New(2),
		isa.Return.New(),
	}, `
	(asm*
		.eval (+ 1 2)
		return)
	`)
}
