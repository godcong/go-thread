package thread

type Threader interface {
	Start()
	Wait()
	Stop()
	PushTo(stepper Stepper, v interface{}) error
	GetThread(stepper Stepper) ThreadRun
	SetThread(stepper Stepper, threader ThreadRun)
	HasThread(stepper Stepper) bool
	SetBaseThread(stepper Stepper, threader Threader)
	IsBase(stepper Stepper) bool
	SetNormalThread(stepper Stepper, threader ThreadRun)
	IsNormal(stepper Stepper) bool
	Register(ops ...Optioner)
}

type threader struct {
}

func NewThreader() Threader
