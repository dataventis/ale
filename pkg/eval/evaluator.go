package eval

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/compiler/procedure"
	"github.com/kode4food/ale/internal/runtime"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/read"
)

// String evaluates the specified raw source
func String(ns env.Namespace, src data.String) (data.Value, error) {
	r := read.FromString(src)
	return Block(ns, r)
}

// Block evaluates a Sequence that a call to FromScanner might produce
func Block(ns env.Namespace, s data.Sequence) (data.Value, error) {
	var res data.Value
	var err error
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		res, err = Value(ns, f)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// Value evaluates the provided Value
func Value(ns env.Namespace, v data.Value) (data.Value, error) {
	defer runtime.NormalizeGoRuntimeErrors()
	e := encoder.NewEncoder(ns)
	if err := generate.Value(e, v); err != nil {
		return nil, err
	}
	e.Emit(isa.Return)
	return encodeAndRun(e)
}

func encodeAndRun(e encoder.Encoder) (data.Value, error) {
	encoded := e.Encode()
	fn, err := procedure.FromEncoded(encoded)
	if err != nil {
		return nil, err
	}
	closure := fn.Call().(data.Procedure)
	return closure.Call(), nil
}
