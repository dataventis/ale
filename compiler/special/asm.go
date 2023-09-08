package special

import (
	"bytes"
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
		labels map[data.Local]isa.Operand
		args   map[data.Local]data.Value
	}

	call struct {
		encoder.Call
		argCount int
	}

	callMap map[data.Local]*call

	toOperandFunc func(data.Value) (isa.Operand, error)
)

// Error messages
const (
	ErrUnknownDirective      = "unknown directive: %s"
	ErrUnexpectedForm        = "unexpected form: %s"
	ErrIncompleteInstruction = "incomplete instruction: %s"
	ErrUnknownLocalType      = "unexpected local type: %s, expected: %s"
	ErrUnexpectedName        = "unexpected local name: %s"
	ErrUnexpectedLabel       = "unexpected label: %s"
)

const (
	MakeEncoder  = data.Local("!make-encoder")
	MakeDetached = data.Local("!make-detached")
	Resolve      = data.Local(".resolve")
	EvalValue    = data.Local(".eval")
	Const        = data.Local(".const")
	Local        = data.Local(".local")
	PushLocals   = data.Local(".push-locals")
	PopLocals    = data.Local(".pop-locals")
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

	cellTypeNames = makeCellTypeNames()
)

// Asm provides indirect access to the Encoder's methods and generators
func Asm(e encoder.Encoder, args ...data.Value) {
	makeAsmEncoder(e).process(data.NewVector(args...))
}

func makeAsmEncoder(e encoder.Encoder) *asmEncoder {
	return &asmEncoder{
		Encoder: e,
		labels:  map[data.Local]isa.Operand{},
		args:    map[data.Local]data.Value{},
	}
}

func (e *asmEncoder) withParams(n data.Locals, v data.Values) *asmEncoder {
	args := make(map[data.Local]data.Value, len(n))
	for i, k := range n {
		args[k] = v[i]
	}
	res := *e
	res.args = args
	return &res
}

func (e *asmEncoder) process(forms data.Sequence) {
	if f, r, ok := forms.Split(); ok {
		if l, ok := f.(data.Local); ok {
			switch l {
			case MakeEncoder:
				e.makeEncoderCall(r)
				return
			case MakeDetached:
				e.makeDetachedEncoder(r)
				return
			}
		}
	}
	e.encode(forms)
}

func (e *asmEncoder) makeDetachedEncoder(forms data.Sequence) {
	ae := makeAsmEncoder(
		encoder.NewEncoder(e.Globals()),
	)
	ae.makeEncoderCall(forms)
}

func (e *asmEncoder) makeEncoderCall(forms data.Sequence) {
	cases := parseParamCases(forms)
	ac := cases.makeArityChecker()
	f := cases.makeFetchers()
	fn := func(e encoder.Encoder, args ...data.Value) {
		if err := ac(len(args)); err != nil {
			panic(err)
		}
		for i, c := range cases {
			if a, ok := f[i](args); ok {
				ae := makeAsmEncoder(e).withParams(c.params, a)
				ae.encode(c.body)
				return
			}
		}
		panic(ErrNoMatchingParamPattern)
	}
	e.Emit(isa.Const, e.AddConstant(encoder.Call(fn)))
}

func (e *asmEncoder) encode(forms data.Sequence) {
	for f, r, ok := forms.Split(); ok; f, r, ok = r.Split() {
		switch name := f.(type) {
		case data.Keyword:
			e.Emit(isa.Label, e.getLabelIndex(name.Name()))
		case data.Local:
			if d, ok := calls[name]; ok {
				if args, rest, ok := take(r, d.argCount); ok {
					d.Call(e, args...)
					r = rest
					continue
				}
				panic(fmt.Errorf(ErrIncompleteInstruction, name))
			}
			panic(fmt.Errorf(ErrUnknownDirective, name))
		default:
			panic(fmt.Errorf(ErrUnexpectedForm, f.String()))
		}
	}
}

func (e *asmEncoder) getLabelIndex(n data.Local) isa.Operand {
	if idx, ok := e.labels[n]; ok {
		return idx
	}
	idx := e.NewLabel()
	e.labels[n] = idx
	return idx
}

func (e *asmEncoder) toOperands(
	oc isa.Opcode, args data.Values,
) []isa.Operand {
	res := make([]isa.Operand, len(args))
	for i, a := range args {
		ao := isa.Effects[oc].Operand
		toOperand := e.getToOperandFor(ao)
		r, err := toOperand(a)
		if err != nil {
			panic(err)
		}
		res[i] = r
	}
	return res
}

func (e *asmEncoder) getToOperandFor(ao isa.ActOn) toOperandFunc {
	switch ao {
	case isa.Locals:
		return e.makeNameToWord()
	case isa.Labels:
		return e.makeLabelToWord()
	default:
		return toOperand
	}
}

func (e *asmEncoder) makeLabelToWord() toOperandFunc {
	return wrapToOperandError(func(val data.Value) (isa.Operand, error) {
		if v, ok := e.resolveEncoderArg(val); ok {
			val = v
		}
		if val, ok := val.(data.Keyword); ok {
			return e.getLabelIndex(val.Name()), nil
		}
		return toOperand(val)
	}, ErrUnexpectedLabel)
}

func (e *asmEncoder) makeNameToWord() toOperandFunc {
	return wrapToOperandError(func(val data.Value) (isa.Operand, error) {
		if v, ok := e.resolveEncoderArg(val); ok {
			val = v
		}
		if val, ok := val.(data.Local); ok {
			if cell, ok := e.ResolveLocal(val); ok {
				return cell.Index, nil
			}
			return 0, fmt.Errorf(ErrUnexpectedName, val)
		}
		return toOperand(val)
	}, ErrUnexpectedName)
}

func (e *asmEncoder) resolveEncoderArg(v data.Value) (data.Value, bool) {
	if v, ok := v.(data.Local); ok {
		res, ok := e.args[v]
		return res, ok
	}
	return nil, false
}

func getInstructionCalls() callMap {
	res := make(callMap, len(isa.Effects))
	for oc, effect := range isa.Effects {
		name := data.Local(strings.CamelToSnake(oc.String()))
		res[name] = func(oc isa.Opcode, ao isa.ActOn) *call {
			return makeEmitCall(oc, ao)
		}(oc, effect.Operand)
	}
	return res
}

func makeEmitCall(oc isa.Opcode, actOn isa.ActOn) *call {
	argCount := 0
	if actOn != isa.Nothing {
		argCount = 1
	}
	return &call{
		Call: func(e encoder.Encoder, args ...data.Value) {
			e.Emit(oc, e.(*asmEncoder).toOperands(oc, args)...)
		},
		argCount: argCount,
	}
}

func getEncoderCalls() callMap {
	return callMap{
		Resolve: {
			Call: func(e encoder.Encoder, args ...data.Value) {
				generate.Symbol(e, args[0].(data.Symbol))
			},
			argCount: 1,
		},
		EvalValue: {
			Call: func(e encoder.Encoder, args ...data.Value) {
				if v, ok := e.(*asmEncoder).resolveEncoderArg(args[0]); ok {
					generate.Value(e, v)
					return
				}
				generate.Value(e, args[0])
			},
			argCount: 1,
		},
		Const: {
			Call: func(e encoder.Encoder, args ...data.Value) {
				if v, ok := e.(*asmEncoder).resolveEncoderArg(args[0]); ok {
					generate.Literal(e, v)
					return
				}
				index := e.AddConstant(args[0])
				e.Emit(isa.Const, index)
			},
			argCount: 1,
		},
		Local: {
			Call: func(e encoder.Encoder, args ...data.Value) {
				name := args[0].(data.Local)
				kwd := args[1].(data.Keyword)
				cellType, ok := cellTypes[kwd]
				if !ok {
					panic(fmt.Errorf(ErrUnknownLocalType, kwd, cellTypeNames))
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

func wrapToOperandError(toOperand toOperandFunc, errStr string) toOperandFunc {
	return func(val data.Value) (isa.Operand, error) {
		res, err := toOperand(val)
		if err != nil {
			return 0, errors.Join(fmt.Errorf(errStr, val), err)
		}
		return res, nil
	}
}

func toOperand(val data.Value) (isa.Operand, error) {
	if val, ok := val.(data.Integer); ok {
		if isa.IsValidOperand(int(val)) {
			return isa.Operand(val), nil
		}
	}
	return 0, fmt.Errorf(isa.ErrExpectedOperand, val)
}

func makeCellTypeNames() string {
	res := mapKeys(cellTypes)
	var buf bytes.Buffer
	for i, s := range res {
		switch {
		case i == len(res)-1:
			buf.WriteString(" or ")
		case i != 0:
			buf.WriteString(" ")
		}
		buf.WriteString(s.String())
	}
	return buf.String()
}

func mapKeys[K comparable, V any](m map[K]V) []K {
	res := make([]K, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}
