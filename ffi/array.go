package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

type arrayWrapper struct {
	typ  reflect.Type
	len  int
	elem Wrapper
}

// Error messages
const (
	ErrValueMustBeSequence = "value must be a sequence"
)

func makeWrappedArray(t reflect.Type) (Wrapper, error) {
	if isMarshaledByteArray(t) {
		return wrapMarshaledByteArray(t)
	}
	w, err := WrapType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &arrayWrapper{
		typ:  t,
		len:  t.Len(),
		elem: w,
	}, nil
}

func (w *arrayWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	vLen := v.Len()
	out := make(data.Values, vLen)
	for i := 0; i < vLen; i++ {
		elem, err := w.elem.Wrap(c, v.Index(i))
		if err != nil {
			return data.Null, err
		}
		out[i] = elem
	}
	return data.NewVector(out...), nil
}

func (w *arrayWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	s, ok := v.(data.Sequence)
	if !ok {
		return _emptyValue, errors.New(ErrValueMustBeSequence)
	}
	in := sequence.ToValues(s)
	out := reflect.New(w.typ).Elem()
	for i, e := range in {
		v, err := w.elem.Unwrap(e)
		if err != nil {
			return _emptyValue, err
		}
		out.Index(i).Set(v)
	}
	return out, nil
}
