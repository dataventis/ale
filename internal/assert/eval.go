package assert

import (
	"gitlab.com/kode4food/ale/bootstrap"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/eval"
	"gitlab.com/kode4food/ale/namespace"
)

var (
	manager = namespace.NewManager()
	ready   bool
)

// Eval will evaluate source code on behalf of the test framework
func (w *Wrapper) Eval(src string) data.Value {
	if !ready {
		bootstrap.Into(manager)
		ready = true
	}
	ns := manager.GetAnonymous()
	return eval.String(ns, data.String(src))
}

// EvalTo will evaluate source code and test for an expected result
func (w *Wrapper) EvalTo(src string, expect data.Value) {
	w.Equal(expect, w.Eval(src))
}

// EvalRaises will evaluate source code and expect an error to be raised
func (w *Wrapper) EvalRaises(src string, err string) {
	defer func() {
		if rec := recover(); rec != nil {
			if v, ok := rec.(data.Value); ok {
				w.String(err, v)
				return
			}
		}
		panic("proper error not raised")
	}()
	w.Eval(src)
}

// PanicWith will evaluate source code and expect a panic to happen
func (w *Wrapper) PanicWith(src string, err error) {
	defer w.ExpectPanic(err.Error())
	w.Eval(src)
}
