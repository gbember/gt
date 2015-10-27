// st_pool_test.go
package pool

import (
	"testing"
)

func TestSTPool(t *testing.T) {
	p := NewSTPool(100)
	st1, _ := p.alloc(1)
	st1.a = 1
	st1.b = "ksgoahqio"
	t.Log(st2)
	t.Log(p.alloc(1))
}
