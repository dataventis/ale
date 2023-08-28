package special

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/strings"
	"github.com/kode4food/ale/runtime/isa"
)

type (
	asmEncoder struct {
		encoder.Encoder
		labels map[data.Name]isa.Index
		args   map[data.Name]data.Value
	}

	call struct {
		encoder.Call
		argCount int
	}

	callMap map[data.Name]*call

	toWordFunc func(data.Value) (isa.Word, error)
)

// Error messages
const (
	ErrUnknownDirective      = "unknown directive: %s"
	ErrUnexpectedForm        = "unexpected form: %s"
	ErrIncompleteInstruction = "incomplete instruction: %s"
	ErrUnknownLocalType      = "unknown local type: %s"
	ErrUnexpectedName        = "unexpected local name: %s"
	ErrUnexpectedLabel       = "unexpected label: %s"
	ErrExpectedWord          = "expected unsigned word: %s"
)

const (
	MakeEncoder = data.Name("!make-encoder")
	Value       = data.Name(".value")
	Const       = data.Name(".const")
	Local       = data.Name(".local")
	PushLocals  = data.Name(".push-locals")
	PopLocals   = data.Name(".pop-locals")
)

var (
	instructionCalls = getInstructionCalls()
	encoderCalls     = getEncoderCalls()
	calls            = mergeCalls(instructionCalls, encoderCalls)

	cellTypes = map[data.Keyword]encoder.CellType{
		data.Keyword("val"):  encoder.ValueCell,
		data.Keyword("ref"):  encoder.ReferenceCell,
		data.Keyword("rest"): encoder.RestCell,
	}
)

// Asm provides indirect access to the Encoder's methods and generators
func Asm(e encoder.Encoder, args ...data.Value) {
	makeEncoder(e).process(data.NewVector(args...))
}

func makeEncoder(e encoder.Encoder) *asmEncoder {
	return &asmEncoder{
		Encoder: e,
		labels:  map[data.Name]isa.Index{},
		args:    map[data.Name]data.Value{},
	}
}

func (e *asmEncoder) withArgs(n data.Names, v data.Values) *asmEncoder {
	res := *e
	args := make(map[data.Name]data.Value, len(n))
	for i, k := range n {
		args[k] = v[i]
	}
	res.args = args
	return &res
}

func (e *asmEncoder) process(forms data.Sequence) {
	if v, r, ok := take(forms, 2); ok {
		if l, ok := v[0].(data.LocalSymbol); ok {
			switch l.Name() {
			case MakeEncoder:
				e.makeEncoder(v[1], r)
				return
			}
		}
	}
	e.encode(forms)
}

func (e *asmEncoder) makeEncoder(arg data.Value, forms data.Sequence) {
	names := parseListArgNames(arg.(data.List))
	fn := func(e encoder.Encoder, args ...data.Value) {
		ae := makeEncoder(e).withArgs(names, args)
		data.AssertFixed(len(names), len(args))
		ae.process(forms)
	}
	e.Emit(isa.Const, e.AddConstant(encoder.Call(fn)))
}

func (e *asmEncoder) encode(forms data.Sequence) {
	for f, r, ok := forms.Split(); ok; f, r, ok = r.Split() {
		switch v := f.(type) {
		case data.Keyword:
			e.Emit(isa.Label, e.getLabelIndex(v))
		case data.LocalSymbol:
			n := v.Name()
			if d, ok := calls[n]; ok {
				if args, rest, ok := take(r, d.argCount); ok {
					d.Call(e, args...)
					r = rest
					continue
				}
				panic(fmt.Errorf(ErrIncompleteInstruction, n))
			}
			panic(fmt.Errorf(ErrUnknownDirective, n))
		default:
			panic(fmt.Errorf(ErrUnexpectedForm, f.String()))
		}
	}
}

func (e *asmEncoder) getLabelIndex(k data.Keyword) isa.Index {
	name := k.Name()
	if idx, ok := e.labels[name]; ok {
		return idx
	}
	idx := e.NewLabel()
	e.labels[name] = idx
	return idx
}

func (e *asmEncoder) toWords(oc isa.Opcode, args data.Values) []isa.Coder {
	res := make([]isa.Coder, len(args))
	for i, a := range args {
		ao := isa.Effects[oc].Operands[i]
		toWord := e.getToWordFor(ao)
		r, err := toWord(a)
		if err != nil {
			panic(err)
		}
		res[i] = r
	}
	return res
}

func (e *asmEncoder) getToWordFor(ao isa.ActOn) toWordFunc {
	switch ao {
	case isa.Locals:
		return e.makeNameToWord()
	case isa.Labels:
		return e.makeLabelToWord()
	default:
		return toWord
	}
}

func (e *asmEncoder) makeLabelToWord() toWordFunc {
	return wrapToWordError(func(val data.Value) (isa.Word, error) {
		if val, ok := val.(data.Keyword); ok {
			return isa.Word(e.getLabelIndex(val)), nil
		}
		return toWord(val)
	}, ErrUnexpectedLabel)
}

func (e *asmEncoder) makeNameToWord() toWordFunc {
	return wrapToWordError(func(val data.Value) (isa.Word, error) {
		if val, ok := val.(data.LocalSymbol); ok {
			if cell, ok := e.ResolveLocal(val.Name()); ok {
				return isa.Word(cell.Index), nil
			}
			return 0, fmt.Errorf(ErrUnexpectedName, val)
		}
		return toWord(val)
	}, ErrUnexpectedName)
}

func getInstructionCalls() callMap {
	res := make(callMap, len(isa.Effects))
	for oc, effect := range isa.Effects {
		name := data.Name(strings.CamelToSnake(oc.String()))
		res[name] = func(oc isa.Opcode, argCount int) *call {
			return makeEmitCall(oc, argCount)
		}(oc, len(effect.Operands))
	}
	return res
}

func makeEmitCall(oc isa.Opcode, argCount int) *call {
	return &call{
		Call: func(e encoder.Encoder, args ...data.Value) {
			e.Emit(oc, e.(*asmEncoder).toWords(oc, args)...)
		},
		argCount: argCount,
	}
}

func getEncoderCalls() callMap {
	return callMap{
		Value: {
			Call: func(e encoder.Encoder, args ...data.Value) {
				if arg, ok := args[0].(data.LocalSymbol); ok {
					if v, ok := e.(*asmEncoder).args[arg.Name()]; ok {
						generate.Value(e, v)
						return
					}
				}
				generate.Value(e, args[0])
			},
			argCount: 1,
		},
		Const: {
			Call: func(e encoder.Encoder, args ...data.Value) {
				index := e.AddConstant(args[0])
				e.Emit(isa.Const, index)
			},
			argCount: 1,
		},
		Local: {
			Call: func(e encoder.Encoder, args ...data.Value) {
				name := args[0].(data.LocalSymbol).Name()
				kwd := args[1].(data.Keyword)
				cellType, ok := cellTypes[kwd]
				if !ok {
					panic(fmt.Errorf(ErrUnknownLocalType, kwd))
				}
				e.AddLocal(name, cellType)
			},
			argCount: 2,
		},
		PushLocals: {
			Call: func(e encoder.Encoder, _ ...data.Value) {
				e.PushLocals()
			},
		},
		PopLocals: {
			Call: func(e encoder.Encoder, _ ...data.Value) {
				e.PopLocals()
			},
		},
	}
}

func mergeCalls(maps ...callMap) callMap {
	res := callMap{}
	for _, m := range maps {
		for k, v := range m {
			if _, ok := res[k]; ok {
				// Programmer error
				panic(fmt.Sprintf("duplicate entry: %s", k))
			}
			res[k] = v
		}
	}
	return res
}

func take(s data.Sequence, count int) (data.Values, data.Sequence, bool) {
	var f data.Value
	var ok bool
	res := make(data.Values, count)
	for i := 0; i < count; i++ {
		if f, s, ok = s.Split(); !ok {
			return nil, nil, false
		}
		res[i] = f
	}
	return res, s, true
}

func wrapToWordError(toWord toWordFunc, errStr string) toWordFunc {
	return func(val data.Value) (isa.Word, error) {
		res, err := toWord(val)
		if err != nil {
			return 0, errors.Join(fmt.Errorf(errStr, val), err)
		}
		return res, nil
	}
}

func toWord(val data.Value) (isa.Word, error) {
	if val, ok := val.(data.Integer); ok {
		if isValidWord(val) {
			return isa.Word(val), nil
		}
	}
	return 0, fmt.Errorf(ErrExpectedWord, val)
}

func isValidWord(i data.Integer) bool {
	return i >= 0 && i <= isa.MaxWord
}
