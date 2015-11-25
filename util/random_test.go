package util

import (
	"math/rand"
	"testing"
)

func BenchmarkMyRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GlobalSyncUniformInt(1, 100)
	}
}

func BenchmarkGoRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Intn(100)
	}
}

type call struct {
	c    chan interface{}
	info interface{}
}
type call_ret struct {
	info interface{}
}

func run() chan interface{} {
	c := make(chan interface{}, 100)
	go func() {
		for {
			i := <-c
			switch i.(type) {
			case *call:
				//TODO 处理call请求
				ret := &call_ret{}
				i.(*call).c <- ret
			case *call_ret:
				//TODO 处理call请求结果
			default:
			}
		}
	}()
	return c
}
