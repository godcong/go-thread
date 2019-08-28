package thread

import (
	"context"

	"go.uber.org/atomic"
	"golang.org/x/xerrors"
)

type Threader interface {
}

// PushFunc ...
type PushFunc func(interface{}) error

// Thread ...
type Thread struct {
	push  PushFunc
	state *atomic.Int32
	done  chan bool
}

// Finished ...
func (t *Thread) Finished() {
	t.SetState(StateStop)
	t.done <- true
}

// Run ...
func (t *Thread) Run(context.Context) {
	panic("implement me")
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
func (t *Thread) BeforeRun(seed Seeder) {
	t.Seeder = seed
}

// AfterRun ...
func (t *Thread) AfterRun(seed Seeder) {
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
