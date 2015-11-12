// mmap_test.go
package mmap

import (
	"testing"
	"time"
)

func TestMMap(t *testing.T) {
	l := &line{point{x: 28, y: 28}, point{x: 308, y: 308}}
	ps := l.getAcossGridNums(5, 1000)
	t.Log(len(ps))
	startTime := time.Now()
	var max int64 = 1000000
	for i := max; i > 0; i-- {
		l.getAcossGridNums(5, 1000)
	}
	td := time.Since(startTime)
	t.Log(td)
	t.Log(td.Nanoseconds() / max)
}
