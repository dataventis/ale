package builtin

import "github.com/kode4food/ale/data"

// Cons adds a value to the beginning of the provided Sequence
var Cons = data.Applicative(func(args ...data.Value) data.Value {
	car := args[0]
	cdr := args[1]
	if p, ok := cdr.(data.Prepender); ok {
		return p.Prepend(car)
	}
	return data.NewCons(car, cdr)
}, 2)

// IsPair returns whether the provided value is a Pair
var IsPair = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.Pair)
	return data.Bool(ok)
}, 1)

// IsCons returns whether the provided value is a Cons cell
var IsCons = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.Cons)
	return data.Bool(ok)
}, 1)
