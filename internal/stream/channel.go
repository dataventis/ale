package stream

import (
	"runtime"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/do"
	"github.com/kode4food/ale/internal/types"
)

type (
	// Emitter is an interface that is used to emit values to a Channel
	Emitter interface {
		Writer
		Closer
	}

	channelEmitter struct {
		ch chan<- data.Value
	}

	channelSequence struct {
		once do.Action
		ch   <-chan data.Value

		result data.Value
		rest   data.Sequence
		ok     bool
	}
)

const (
	// EmitKey is the key used to emit to a Channel
	EmitKey = data.Keyword("emit")

	// SequenceKey is the key used to retrieve the Sequence from a Channel
	SequenceKey = data.Keyword("seq")
)

var (
	emptyResult = data.Null

	channelSequenceType = types.MakeBasic("channel-sequence")
)

// NewChannel produces an Emitter and Sequence pair
func NewChannel(size int) *data.Object {
	ch := make(chan data.Value, size)
	e := NewChannelEmitter(ch)
	s := NewChannelSequence(ch)
	return data.NewObject(
		data.NewCons(EmitKey, bindWriter(e)),
		data.NewCons(CloseKey, bindCloser(e)),
		data.NewCons(SequenceKey, s),
	)
}

// NewChannelEmitter produces an Emitter for sending values to a Go chan
func NewChannelEmitter(ch chan<- data.Value) Emitter {
	r := &channelEmitter{
		ch: ch,
	}
	runtime.SetFinalizer(r, func(e *channelEmitter) {
		defer func() { recover() }()
		close(ch)
	})
	return r
}

// Write will send a Value to the Go chan
func (e *channelEmitter) Write(v data.Value) {
	e.ch <- v
}

// Close will Close the Go chan
func (e *channelEmitter) Close() {
	runtime.SetFinalizer(e, nil)
	close(e.ch)
}

// NewChannelSequence produces a new Sequence whose values come from a Go chan
func NewChannelSequence(ch <-chan data.Value) data.Sequence {
	return &channelSequence{
		once:   do.Once(),
		ch:     ch,
		result: emptyResult,
		rest:   data.Null,
	}
}

func (c *channelSequence) resolve() *channelSequence {
	c.once(func() {
		result, ok := <-c.ch
		if !ok {
			return
		}
		c.ok = ok
		c.result = result
		c.rest = NewChannelSequence(c.ch)
	})

	return c
}

func (c *channelSequence) IsEmpty() bool {
	return !c.resolve().ok
}

func (c *channelSequence) Car() data.Value {
	return c.resolve().result
}

func (c *channelSequence) Cdr() data.Value {
	return c.resolve().rest
}

func (c *channelSequence) Split() (data.Value, data.Sequence, bool) {
	r := c.resolve()
	return r.result, r.rest, r.ok
}

func (c *channelSequence) Prepend(v data.Value) data.Sequence {
	return &channelSequence{
		once:   do.Never(),
		ok:     true,
		result: v,
		rest:   c,
	}
}

func (c *channelSequence) Type() types.Type {
	return channelSequenceType
}

func (c *channelSequence) Equal(v data.Value) bool {
	return c == v
}

func (c *channelSequence) String() string {
	return data.DumpString(c)
}
