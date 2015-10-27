// queue_test.go
package pool

import (
	"testing"
)

type tt struct {
	a int
	b int
	c int
	d int
}

func TestQueue(t *testing.T) {
	queue := NewQueue(10)
	t.Logf("=====%p", &queue.bs[0])
	t.Logf("=====%p", &queue.bs[1])
	queue.Push(tt{a: 10})
	t.Logf("=====%v", queue.bs[0])
}
