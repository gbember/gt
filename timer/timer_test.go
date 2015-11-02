// timer_test.go
package timer

import (
	"log"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := NewTimer()
	timer.AddFun(10*time.Second, func() { log.Println("111111111") })
	timer.AddFun(5*time.Second, func() { log.Println("222222222") })
	timer.AddFun(2*time.Second, func() { log.Println("3333333333") })
	<-timer.T.C
	timer.Run()
	<-timer.T.C
	timer.Run()
	<-timer.T.C
	timer.Run()
}
