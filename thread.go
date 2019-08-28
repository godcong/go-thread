package thread

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/godcong/go-trait"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

var log = trait.NewZapSugar()

// DefaultInterval ...
const DefaultInterval = 30 * time.Second

// Runner ...
type Runner interface {
	Runnable
	BeforeRun(manager Manager)
	AfterRun(manager Manager)
}

// Basic ...
type Basic interface {
	SetInterval(duration time.Duration)
	Interval() time.Duration
	ID() string
	Manager() Manager
}

// NoStatusThreader ...
type NoStatusThreader interface {
	Threader
	State() State
	SetState(state State)
	Done() <-chan bool
	Finished()
}

// Threader ...
type Threader interface {
	Runner
	Basic
	Pusher
}

// CallAble ...
type CallAble interface {
	Call(Threader, interface{}) error
}

// Runnable ...
type Runnable interface {
	Run(context.Context)
}

// Pusher ...
type Pusher interface {
	Push(interface{}) error
}

// State ...
type State int

// State ...
const (
	StateWaiting State = iota
	StateRunning
	StateStop
)

// PushFunc ...
type PushFunc func(interface{})

// Thread ...
type Thread struct {
	manager  Manager
	id       string
	interval time.Duration
	state    *int32
	done     chan bool
	cb       chan interface{}
	CallAble
}

// Manager ...
func (t *Thread) Manager() Manager {
	return t.manager
}

// SetInterval ...
func (t *Thread) SetInterval(duration time.Duration) {
	t.interval = duration
}

// Interval ...
func (t *Thread) Interval() time.Duration {
	return t.interval
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
			e := t.Call(t, cb)
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
	if v != nil {
		go func(cb chan<- interface{}, v interface{}) {
			cb <- v
		}(t.cb, v)
		return nil
	}
	return xerrors.New("null push function")
}

// BeforeRun ...
func (t *Thread) BeforeRun(manager Manager) {
	t.manager = manager
}

// AfterRun ...
func (t *Thread) AfterRun(manager Manager) {
}

// State ...
func (t *Thread) State() State {
	return State(atomic.LoadInt32(t.state))
}

// ID ...
func (t *Thread) ID() string {
	return t.id
}

// Done ...
func (t *Thread) Done() <-chan bool {
	return t.done
}

// NewThreader ...
func NewThreader(call CallAble) Threader {
	state := int32(StateWaiting)
	return &Thread{
		id:       uuid.New().String(),
		interval: DefaultInterval,
		state:    &state,
		done:     make(chan bool),
		CallAble: call,
	}
}
