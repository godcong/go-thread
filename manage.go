package thread

import (
	"context"
	"errors"
	"time"

	"go.uber.org/atomic"
)

// Options ...
type Options func(manage *Manage)

// Manager ...
type Manager interface {
	Start()
	Wait()
	Stop()
	PushTo(id string, v interface{}) error
	GetThread(id string) Threader
	AddThread(threader Threader) string
	HasThread(id string) bool
	Register(ops ...Options)
}

// Manage ...
type Manage struct {
	ctx       context.Context
	cancel    context.CancelFunc
	threaders map[string]Threader
}

// Start ...
func (m *Manage) Start() {

}

// Wait ...
func (m *Manage) Wait() {
	state := 1
	for state > 0 {
		for _, t := range m.threaders {
			if StateWaiting != t.State() {
				state = 2
				time.Sleep(15 * time.Second)
				break
			}
		}
		if state == 2 {
			state = 1
		} else {
			state = 0
		}
	}
	log.Info("base done")
	m.Done()
}

// Done ...
func (m *Manage) Done() {
	count := atomic.NewInt32(0)
	for i := range m.threaders {
		go func(threader Threader) {
			//if base.State() != StateStop {
			<-threader.Done()
			//}
			count.Add(1)
		}(m.threaders[i])
		m.threaders[i].Finished()
	}
	for {
		if count.Load() == int32(len(m.threaders)) {
			return
		}
	}
}

// Stop ...
func (m *Manage) Stop() {
	if m.cancel != nil {
		m.cancel()
	}
}

// PushTo ...
func (m *Manage) PushTo(id string, v interface{}) error {
	if threader, b := m.threaders[id]; b {
		return threader.Push(v)
	}
	return errors.New("thread not found")
}

// GetThread ...
func (m *Manage) GetThread(id string) Threader {
	return m.threaders[id]
}

// SetThread ...
func (m *Manage) AddThread(threader interface{}) {

}

// HasThread ...
func (m *Manage) HasThread(id string) bool {
	_, b := m.threaders[id]
	return b
}

// Register ...
func (m *Manage) Register(ops ...Options) {
	for _, v := range ops {
		v(m)
	}
}

// NewManager create a thread manager
// manager is not always thread safe
// so all register need before running
func NewManager(ts ...Threader) Manager {
	return nil
}
