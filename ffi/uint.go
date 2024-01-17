package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/debug"
)

type (
	uintWrapper   reflect.Kind
	uint64Wrapper reflect.Kind
	uint32Wrapper reflect.Kind
	uint16Wrapper reflect.Kind
	uint8Wrapper  reflect.Kind
)

var (
	uintZero   = reflect.ValueOf(uint(0))
	uint64zero = reflect.ValueOf(uint64(0))
	uint32zero = reflect.ValueOf(uint32(0))
	uint16zero = reflect.ValueOf(uint16(0))
	uint8zero  = reflect.ValueOf(uint8(0))
)

func makeWrappedUnsignedInt(t reflect.Type) (Wrapper, error) {
	k := t.Kind()
	switch k {
	case reflect.Uint:
		return uintWrapper(k), nil
	case reflect.Uint64:
		return uint64Wrapper(k), nil
	case reflect.Uint32:
		return uint32Wrapper(k), nil
	case reflect.Uint16:
		return uint16Wrapper(k), nil
	case reflect.Uint8:
		return uint8Wrapper(k), nil
	default:
		panic(debug.ProgrammerError("uint kind is incorrect"))
	}
}

func (uintWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Uint()), nil
}

func (uintWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if i, ok := v.(data.Integer); ok {
		return reflect.ValueOf(uint(i)), nil
	}
	return uintZero, errors.New(ErrValueMustBeInteger)
}

func (uint64Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Uint()), nil
}

func (uint64Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if i, ok := v.(data.Integer); ok {
		return reflect.ValueOf(uint64(i)), nil
	}
	return uint64zero, errors.New(ErrValueMustBeInteger)
}

func (uint32Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Uint()), nil
}

func (uint32Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if i, ok := v.(data.Integer); ok {
		return reflect.ValueOf(uint32(i)), nil
	}
	return uint32zero, errors.New(ErrValueMustBeInteger)
}

func (uint16Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Uint()), nil
}

func (uint16Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if i, ok := v.(data.Integer); ok {
		return reflect.ValueOf(uint16(i)), nil
	}
	return uint16zero, errors.New(ErrValueMustBeInteger)
}

func (uint8Wrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	return data.Integer(v.Uint()), nil
}

func (uint8Wrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if i, ok := v.(data.Integer); ok {
		return reflect.ValueOf(uint8(i)), nil
	}
	return uint8zero, errors.New(ErrValueMustBeInteger)
}
