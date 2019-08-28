package thread

type Manager interface {
	Start()
	Wait()
	Stop()
	PushTo(id string, v interface{}) error
	GetThread(id string) ThreadRun
	SetThread(id string, threader ThreadRun)
	HasThread(id string) bool
	Register(ops ...Optionor)
}

type manage struct {
}

func NewManager(t ...Threader) Manager {

}
