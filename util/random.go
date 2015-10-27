// random.go
//随机值
package util

import (
	"sync"
	"time"
)

const (
	PRIME1 = 30269
	PRIME2 = 30307
	PRIME3 = 30323

	F_PRIME1 float32 = 30269
	F_PRIME2 float32 = 30307
	F_PRIME3 float32 = 30323
)

type _rand struct {
	_mut   sync.Mutex
	prime1 int
	prime2 int
	prime3 int
}

var (
	_globalRand *_rand = Seed()
	_globalSeed int
	_globalMut  sync.Mutex
)

//构造一个随机器
func Seed() *_rand {
	_globalMut.Lock()
	addp := _globalSeed
	_globalSeed++
	_globalMut.Unlock()

	ms := time.Now().Nanosecond() / 1000000
	prime1 := ms / 1000000000
	ms = ms - prime1
	prime2 := ms/1000 + addp + 1
	prime3 := (ms-prime2)*1000 + addp + 2
	return &_rand{prime1: prime1, prime2: prime2, prime3: prime3}
}

//设置随机器的随机基数
func (_r *_rand) seed() (int, int, int) {
	prime1 := (_r.prime1 * 171) % PRIME1
	prime2 := (_r.prime2 * 172) % PRIME2
	prime3 := (_r.prime3 * 170) % PRIME3
	_r.prime1 = prime1
	_r.prime2 = prime2
	_r.prime3 = prime3
	return prime1, prime2, prime3
}

//全局并发不安全的随机0-1之间的小数(包括0不包括1)
func GlobalUniformFloat() float32 {
	return _globalRand.UniformFloat()
}

//全局并发安全的随机0-1之间的小数(包括0不包括1)
func GlobalSyncUniformFloat() float32 {
	return _globalRand.SyncUniformFloat()
}

//全局并发不安全的随机1-max之间的小数(包括1和max)
func GlobalUniformInt(max int) int {
	return _globalRand.UniformInt(max)
}

//全局并发不安全的随机1-max之间的小数(包括1和max)
func GlobalSyncUniformInt(max int) int {
	return _globalRand.SyncUniformInt(max)
}

//并发不安全的随机0-1之间的小数(包括0不包括1)
func (_r *_rand) UniformFloat() float32 {
	prime1, prime2, prime3 := _r.seed()
	r := float32(prime1)/F_PRIME1 + float32(prime2)/F_PRIME2 + float32(prime3)/F_PRIME3
	return r - float32(int(r))
}

//并发安全的随机0-1之间的小数(包括0不包括1)
func (_r *_rand) SyncUniformFloat() float32 {
	_r._mut.Lock()
	r := _r.UniformFloat()
	_r._mut.Unlock()
	return r
}

//并发不安全的随机1-max之间的小数(包括1和max)
func (_r *_rand) UniformInt(max int) int {
	if max < 1 {
		panic("random UniformInt param max must >= 1")
	}
	f := _r.UniformFloat()
	return int(f*float32(max)) + 1
}

//并发不安全的随机1-max之间的小数(包括1和max)
func (_r *_rand) SyncUniformInt(max int) int {
	_r._mut.Lock()
	r := _r.UniformInt(max)
	_r._mut.Unlock()
	return r
}
