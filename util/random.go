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

type Rand interface {
	//并发不安全的随机0-1之间的小数(包括0不包括1)
	UniformFloat() float32
	//并发安全的随机0-1之间的小数(包括0不包括1)
	SyncUniformFloat() float32
	//并发不安全的随机min-max之间的小数(包括min和max)
	UniformInt(min int, max int) int
	//并发不安全的随机max-max之间的小数(包括min和max)
	SyncUniformInt(min int, max int) int
}

type _rand struct {
	_mut   sync.Mutex
	prime1 int
	prime2 int
	prime3 int
}

var (
	_globalRand Rand = Seed()
	_globalSeed int
	_globalMut  sync.Mutex
)

//构造一个随机器
func Seed() Rand {
	return SeedInt(0)
}

//构造一个随机器(带个随机值)
func SeedInt(s int) Rand {
	prime1 := time.Now().Nanosecond()
	prime2 := prime1/1000 + s + 1
	prime3 := (prime1-prime2)*1000 + s + 2
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

//全局并发不安全的随机min-max之间的整数(包括1和max)
func GlobalUniformInt(min int, max int) int {
	return _globalRand.UniformInt(min, max)
}

//全局并发不安全的随机1-max之间的整数(包括1和max)
func GlobalSyncUniformInt(min int, max int) int {
	return _globalRand.SyncUniformInt(min, max)
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

//并发不安全的随机min-max之间的小数(包括1和max)
func (_r *_rand) UniformInt(min int, max int) int {
	if max < min {
		panic("random UniformInt param max must >= 1")
	}
	f := _r.UniformFloat()
	return int(f*float32(max-min)) + min
}

//并发不安全的随机max-max之间的小数(包括1和max)
func (_r *_rand) SyncUniformInt(min int, max int) int {
	_r._mut.Lock()
	r := _r.UniformInt(min, max)
	_r._mut.Unlock()
	return r
}
