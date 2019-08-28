package thread

import (
	"context"
	"time"

	"go.uber.org/atomic"
	"golang.org/x/xerrors"
)

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

// PushFunc ...
type PushFunc func(interface{}) error

// Thread ...
type Thread struct {
	Threader
	interval time.Duration
	push     PushFunc
	state    *atomic.Int32
	done     chan bool
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
			e := cb.Call(m)
			if e != nil {
				log.Error(e)
			}
		case <-time.After(t.interval):
			log.Info("info time out")
			m.SetState(StateWaiting)
		}
	}
	close(m.cb)
	m.Finished()
}

// SetState ...
func (t *Thread) SetState(state State) {
	t.state.Store(int32(state))
}

// Push ...
func (t *Thread) Push(v interface{}) error {
	if t.push != nil {
		return t.push(v)
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
	return State(t.state.Load())
}

// Done ...
func (t *Thread) Done() <-chan bool {
	return t.done
}

// NewThread ...
func NewThread() *Thread {
	return &Thread{
		state: atomic.NewInt32(int32(StateWaiting)),
		done:  make(chan bool),
	}
}
