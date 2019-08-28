package thread

// Options ...
type Options func(manage *Manage)

// Manager ...
type Manager interface {
	Start()
	Wait()
	Stop()
	PushTo(id string, v interface{}) error
	GetThread(id string) Runner
	SetThread(id string, runner Runner)
	HasThread(id string) bool
	Register(ops ...Options)
}

// Manage ...
type Manage struct {
	threaders map[string]Threader
}

// Start ...
func (m *Manage) Start() {
	panic("implement me")
}

// Wait ...
func (m *Manage) Wait() {
	panic("implement me")
}

// Stop ...
func (m *Manage) Stop() {
	panic("implement me")
}

// PushTo ...
func (m *Manage) PushTo(id string, v interface{}) error {
	panic("implement me")
}

// GetThread ...
func (m *Manage) GetThread(id string) Runner {
	panic("implement me")
}

// SetThread ...
func (m *Manage) SetThread(id string, runner Runner) {
	panic("implement me")
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
	m := &Manage{
		threaders: make(map[string]Threader),
	}
	for i, t := range ts {
		m.threaders[t.ID()] = ts[i]
	}
	return m
}
