package thread

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Thread struct {
	interval time.Duration
	id       string
	ctx      context.Context
	cancel   context.CancelFunc
	run      Runnable
}

func (thread *Thread) Start() {
	thread.ctx, thread.cancel = context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-thread.ctx.Done():
				return
			default:
				e := thread.run()
				if e != nil {
					panic(e)
				}
			}
			time.Sleep(thread.interval)
		}
	}()
}

func (thread *Thread) Stop() {
	if thread.cancel != nil {
		thread.cancel()
	}
}

type Runnable func() error

func New(run Runnable) *Thread {
	if run == nil {
		panic("must input an runnable function")
	}
	return &Thread{
		interval: 300,
		id:       uuid.New().String(),
		run:      run,
	}
}
