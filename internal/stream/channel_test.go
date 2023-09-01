package stream_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/stream"
)

func TestChannel(t *testing.T) {
	as := assert.New(t)

	e, seq := stream.NewChannel(0)
	seq = seq.(data.Prepender).Prepend(F(1))
	as.Contains(":type channel-emitter", e)
	as.Contains(":type channel-sequence", seq)

	var wg sync.WaitGroup

	gen := func() {
		e.Write(F(2))
		time.Sleep(time.Millisecond * 50)
		e.Write(F(3))
		time.Sleep(time.Millisecond * 30)
		e.Write(S("foo"))
		time.Sleep(time.Millisecond * 10)
		e.Write(S("bar"))
		e.Close()
		wg.Done()
	}

	check := func() {
		f, _, ok := seq.Split()
		as.Number(1, f)
		as.True(ok)

		as.Number(1, seq.Car())
		as.Number(2, seq.Cdr().(data.Pair).Car())
		as.Number(3, seq.Cdr().(data.Pair).Cdr().(data.Pair).Car())
		as.False(seq.Cdr().(data.Pair).Cdr().(data.Pair).Cdr().(data.Sequence).
			IsEmpty())
		as.String("foo", seq.Cdr().(data.Pair).Cdr().(data.Pair).
			Cdr().(data.Pair).Car())
		as.False(seq.Cdr().(data.Pair).Cdr().(data.Pair).Cdr().(data.Pair).
			Cdr().(data.Sequence).IsEmpty())
		as.String("bar", seq.Cdr().(data.Pair).Cdr().(data.Pair).
			Cdr().(data.Pair).Cdr().(data.Pair).Car())
		as.True(seq.Cdr().(data.Pair).Cdr().(data.Pair).Cdr().(data.Pair).
			Cdr().(data.Pair).Cdr().(data.Sequence).IsEmpty())
		wg.Done()
	}

	wg.Add(4)
	go check()
	go check()
	go gen()
	go check()
	wg.Wait()
}

func TestChannelError(t *testing.T) {
	as := assert.New(t)

	e, seq := stream.NewChannel(2)
	e.Write(S("hello"))
	e.Error(fmt.Errorf("boom"))

	f, r, ok := seq.Split()
	as.True(ok)
	as.Equal(S("hello"), f)
	as.NotNil(r)

	defer as.ExpectPanic("boom")
	_ = r.Car()
}
