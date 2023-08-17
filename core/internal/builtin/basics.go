package builtin

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/read"
	"github.com/kode4food/ale/runtime"
)

// Raise will cause a panic
var Raise = data.Applicative(func(args ...data.Value) data.Value {
	panic(args[0])
}, 1)

// Recover invokes a function and runs a recovery function if Go panics
var Recover = data.Applicative(func(args ...data.Value) (res data.Value) {
	body := args[0].(data.Function)
	rescue := args[1].(data.Function)

	defer func() {
		if rec := recover(); rec != nil {
			switch rec := runtime.NormalizeGoRuntimeError(rec).(type) {
			case data.Value:
				res = rescue.Call(rec)
			case error:
				res = rescue.Call(data.String(rec.Error()))
			default:
				// Programmer error
				panic("rescue returned an invalid result")
			}
		}
	}()

	return body.Call()
}, 2)

// Defer invokes a cleanup function, no matter what has happened
var Defer = data.Applicative(func(args ...data.Value) (res data.Value) {
	body := args[0].(data.Function)
	cleanup := args[1].(data.Function)

	defer cleanup.Call()
	return body.Call()
}, 2)

// Read performs the standard LISP read of a string
var Read = data.Applicative(func(args ...data.Value) data.Value {
	v := args[0]
	s := v.(data.Sequence)
	if v, ok := data.Last(read.FromString(sequence.ToStr(s))); ok {
		return v
	}
	return data.Nil
}, 1)
