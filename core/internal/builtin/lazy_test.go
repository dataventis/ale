package builtin_test

import (
	"testing"

	"github.com/kode4food/ale/core/internal/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestLazySequence(t *testing.T) {
	as := assert.New(t)

	var i int
	var fn data.Function

	fn = data.Applicative(func(...data.Value) data.Value {
		if i < 10 {
			res := builtin.Cons.Call(
				data.Integer(i),
				builtin.LazySequence.Call(fn),
			)
			i++
			return res
		}
		return data.Nil
	}, 0)

	s := builtin.LazySequence.Call(fn).(data.Sequence)
	as.String(`(0 1 2 3 4 5 6 7 8 9)`, data.MakeSequenceStr(s))
}

func TestRangeEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(reduce
			(lambda (x y) (+ x y))
			(range 1 5 1))
	`, F(10))

	as.EvalTo(`
		(reduce
			(lambda (x y) (+ x y))
			(range 5 1 -1))
	`, F(14))
}

func TestMapAndFilterEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(reduce
			(lambda (x y) (+ x y))
			(map
				(lambda (x) (* x 2))
				(filter
					(lambda (x) (<= x 5))
					[1 2 3 4 5 6 7 8 9 10])))
	`, F(30))
}

func TestMapParallelEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(seq->vector
			(map +
				[1 2 3 4]
				'(2 4 6 8)
				(range 20 30)))
	`, S("[23 27 31 35]"))
}

func TestReduceEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define x '(1 2 3 4))
		(reduce + x)
	`, F(10))

	as.EvalTo(`
		(define y (concat '(1 2 3 4) [5 6 7 8]))
		(reduce + y)
	`, F(36))

	as.EvalTo(`
		(define y (concat '(1 2 3 4) [5 6 7 8]))
		(reduce + 10 y)
	`, F(46))
}

func TestTakeDropEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define x (concat '(1 2 3 4) [5 6 7 8]))
		(nth (apply vector (take 6 x)) 5)
	`, F(6))

	as.EvalTo(`
		(define x (concat '(1 2 3 4) [5 6 7 8]))
		(nth (apply vector (drop 3 x)) 0)
	`, F(4))

	err := unexpectedTypeError("integer", "sequence")
	as.PanicWith(`(last! (drop 99 57))`, err)
	as.PanicWith(`(last! (take 99 57))`, err)
}

func TestLazySeqEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(reduce
			(lambda (x y) (+ x y))
			(lazy-seq (cons 1 (lazy-seq [2 3]))))
	`, F(6))

	as.EvalTo(`
		(length (seq->vector (lazy-seq '())))
	`, F(0))
}

func TestForEachLoopEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(let* ([ch (chan)]
			   [emit (:emit ch)]
			   [close (:close ch)]
			   [seq (:seq ch)])
			(go
				(for-each ([i (range 1 5 1)]
				           [j (range 1 10 2)])
					(emit (* i j)))
				(close))
			(reduce (lambda (x y) (+ x y)) seq))
	`, F(250))
}
