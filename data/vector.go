package data

import (
	"bytes"
	"math/rand"

	"github.com/kode4food/ale/types"
)

type (
	// Vector is a fixed-length array of Values
	Vector interface {
		vector() // marker
		Sequence
		ValuerSequence
		Prepender
		Appender
		Reverser
		RandomAccess
		Caller
		Valuer
	}

	vector Values
)

// EmptyVector represents an empty Vector
var (
	EmptyVector = vector{}

	vectorHash = rand.Uint64()
)

// NewVector creates a new Vector instance
func NewVector(v ...Value) Vector {
	return vector(v)
}

func (vector) vector() {}

func (v vector) Values() Values {
	return Values(v)
}

func (v vector) Count() int {
	return len(v)
}

func (v vector) ElementAt(index int) (Value, bool) {
	if index >= 0 && index < len(v) {
		return v[index], true
	}
	return Nil, false
}

func (v vector) First() Value {
	if len(v) > 0 {
		return v[0]
	}
	return Nil
}

func (v vector) Rest() Sequence {
	if len(v) > 1 {
		return v[1:]
	}
	return EmptyVector
}

func (v vector) IsEmpty() bool {
	return len(v) == 0
}

func (v vector) Split() (Value, Sequence, bool) {
	lv := len(v)
	if lv > 1 {
		return v[0], v[1:], true
	} else if lv == 1 {
		return v[0], EmptyVector, true
	}
	return Nil, EmptyVector, false
}

func (v vector) Car() Value {
	return v.First()
}

func (v vector) Cdr() Value {
	return v.Rest()
}

func (v vector) Prepend(e Value) Sequence {
	return append(vector{e}, v...)
}

func (v vector) Append(e Value) Sequence {
	return append(v, e)
}

func (v vector) Reverse() Sequence {
	vl := len(v)
	if vl <= 1 {
		return v
	}
	res := make(vector, vl)
	for i, j := 0, vl-1; j >= 0; i, j = i+1, j-1 {
		res[i] = v[j]
	}
	return res
}

func (v vector) Call(args ...Value) Value {
	return indexedCall(v, args)
}

func (v vector) Convention() Convention {
	return ApplicativeCall
}

func (v vector) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

func (v vector) Equal(r Value) bool {
	if r, ok := r.(vector); ok {
		if len(v) != len(r) {
			return false
		}
		for i, elem := range r {
			if !v[i].Equal(elem) {
				return false
			}
		}
		return true
	}
	return false
}

func (v vector) String() string {
	var b bytes.Buffer
	b.WriteString("[")
	for i, e := range v {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(MaybeQuoteString(e))
	}
	b.WriteString("]")
	return b.String()
}

func (vector) Type() types.Type {
	return types.AnyVector
}

func (v vector) HashCode() uint64 {
	h := vectorHash
	for _, e := range v {
		h *= HashCode(e)
	}
	return h
}
