package generate

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

// Symbol encodes a symbol retrieval
func Symbol(e encoder.Encoder, s data.Symbol) error {
	if l, ok := s.(data.Local); ok {
		_, err := resolveLocal(e, l)
		return err
	}
	return resolveGlobal(e, s)
}

// ReferenceSymbol encodes a potential symbol retrieval and dereference
func ReferenceSymbol(e encoder.Encoder, s data.Symbol) error {
	switch s := s.(type) {
	case data.Local:
		c, err := resolveLocal(e, s)
		if err != nil {
			return err
		}
		if c != nil && c.Type == encoder.ReferenceCell {
			e.Emit(isa.Deref)
		}
	default:
		return resolveGlobal(e, s)
	}
	return nil
}

func resolveLocal(
	e encoder.Encoder, l data.Local,
) (*encoder.ScopedCell, error) {
	if s, ok := e.ResolveScoped(l); ok {
		switch s.Scope {
		case encoder.LocalScope:
			c, _ := e.ResolveLocal(l)
			e.Emit(isa.Load, c.Index)
		case encoder.ArgScope:
			c, _ := e.ResolveParam(l)
			if c.Type == encoder.RestCell {
				e.Emit(isa.RestArg, c.Index)
			} else {
				e.Emit(isa.Arg, c.Index)
			}
		case encoder.ClosureScope:
			c, _ := e.ResolveClosure(l)
			e.Emit(isa.Closure, c.Index)
		default:
			panic(debug.ProgrammerError("unknown scope type"))
		}
		return s, nil
	}
	return nil, resolveGlobal(e, l)
}

func resolveGlobal(e encoder.Encoder, s data.Symbol) error {
	entry, _, err := env.ResolveSymbol(e.Globals(), s)
	if err != nil {
		return err
	}
	if entry.IsBound() {
		v, _ := entry.Value()
		return Literal(e, v)
	}
	if err := Literal(e, s); err != nil {
		return err
	}
	e.Emit(isa.Resolve)
	return nil
}
