package thread

import (
	"context"
	"time"

	"sync/atomic"

	"github.com/godcong/go-trait"
	"golang.org/x/xerrors"
)

var log = trait.NewZapSugar()

// State ...
type State int

// State ...
const (
	StateWaiting State = iota
	StateRunning
	StateStop
)

type Threader interface {
}

type CallAble interface {
	Call(*Thread, interface{}) error
}

// PushFunc ...
type PushFunc func(interface{})

// Thread ...
type Thread struct {
	Threader
	interval time.Duration
	push     PushFunc
	state    *int32
	done     chan bool
	cb       chan CallAble
	call     CallAble
}

// Finished ...
func (t *Thread) Finished() {
	t.SetState(StateStop)
	t.done <- true
}

// Run ...
func (t *Thread) Run(ctx context.Context) {
ThreadEnd:
	for {
		select {
		case <-ctx.Done():
			break ThreadEnd
		case cb := <-t.cb:
			if cb == nil {
				break ThreadEnd
			}
			t.SetState(StateRunning)
			e := cb.Call(t)
			if e != nil {
				log.Error(e)
			}
		case <-time.After(t.interval):
			log.Info("info time out")
			t.SetState(StateWaiting)
		}
	}
	close(t.cb)
	t.Finished()
}

// SetState ...
func (t *Thread) SetState(state State) {
	atomic.StoreInt32(t.state, int32(state))
}

// Push ...
func (t *Thread) Push(v interface{}) error {
	if t.push != nil {
		go func(p PushFunc, v interface{}) {
			p(v)
		}(t.push, v)
		return nil
	}
	return xerrors.New("null push function")
}

// BeforeRun ...
func (t *Thread) BeforeRun(thread Threader) {
	t.Threader = thread
}

// AfterRun ...
func (t *Thread) AfterRun(thread Threader) {
}

// State ...
func (t *Thread) State() State {
	return State(atomic.LoadInt32(t.state))
}

// Done ...
func (t *Thread) Done() <-chan bool {
	return t.done
}

// NewThread ...
func NewThread(call CallAble) *Thread {
	state := int32(StateWaiting)
	return &Thread{
		state: (*int32)(&state),
		done:  make(chan bool),
		call:  call,
	}
}
