package timer

import (
	"container/heap"
	"time"
)

const (
	_DEFAULT_TIMER_D = time.Hour * 24
)

type timerFun struct {
	d   int64
	fun func()
}

type timerFunList []*timerFun

func (this timerFunList) Len() int           { return len(this) }
func (this timerFunList) Less(i, j int) bool { return this[i].d < this[j].d }
func (this timerFunList) Swap(i, j int)      { this[i], this[j] = this[j], this[i] }
func (this *timerFunList) Push(x interface{}) {
	*this = append(*this, x.(*timerFun))
}
func (this *timerFunList) Pop() interface{} {
	old := *this
	n := len(old)
	x := old[n-1]
	*this = old[0 : n-1]
	return x
}

type TimerFunQueue struct {
	*time.Timer
	d   int64
	tfl timerFunList
}

func NewTimerFunQueue() *TimerFunQueue {
	tfq := new(TimerFunQueue)
	heap.Init(&tfq.tfl)
	tfq.defaultSet()
	return tfq
}

func (tfq *TimerFunQueue) defaultSet() {
	tfq.d = 0
	tfq.tfl = tfq.tfl[0:]
	tfq.Timer = time.NewTimer(_DEFAULT_TIMER_D)
}

func (tfq *TimerFunQueue) AddAfterTimerFun(d time.Duration, fun func()) {
	if d <= 0 {
		fun()
		return
	}
	rd := time.Now().Add(d).UnixNano()
	heap.Push(&tfq.tfl, &timerFun{d: rd, fun: fun})
	if tfq.d == 0 || rd < tfq.d {
		tfq.d = rd
		tfq.Reset(d)
	}

}

func (tfq *TimerFunQueue) Run() {
	for tfq.tfl.Len() > 0 {
		tf := heap.Pop(&tfq.tfl).(*timerFun)
		if tf.d > tfq.d {
			tfq.d = 0
			tfq.AddAfterTimerFun(time.Duration(tf.d-time.Now().UnixNano()), tf.fun)
			return
		}
		tf.fun()
	}
	tfq.defaultSet()
}
