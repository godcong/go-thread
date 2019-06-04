package thread

import (
	"fmt"
	"testing"
	"time"
)

func TestThread_Start(t *testing.T) {
	th := New(func() error {
		fmt.Println(time.Now())
		time.Sleep(10 * time.Second)
		return nil
	})
	th.Start()

	time.Sleep(11 * time.Second)
	th.Stop()

}
