// synccond_test.go
package concurrent

import (
	"runtime"
	"testing"
	"time"
)

func TestConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmap := NewConcurrentMap(900000)
	go func() {
		for {
			cmap.DeleteAnyOneWait()
		}
	}()
	go func() {
		for {
			cmap.DeleteAnyOneWait()
		}
	}()
	go func() {
		for {
			cmap.DeleteAnyOneWait()
		}
	}()
	time.Sleep(time.Second)
	var max int64 = 900000
	startTime := time.Now()
	for i := max; i > 0; i-- {
		cmap.PutAndSignal(i, i)
	}
	t.Log(time.Since(startTime).Nanoseconds() / max)
}
