// timer.go
package timer

import "time"

type Timer struct {
	T     *time.Timer
	tf    *timerFun
	tfMap map[*timerFun]bool
}

type timerFun struct {
	d   int64
	fun func()
}

func NewTimer() *Timer {
	timer := new(Timer)
	timer.T = &time.Timer{C: make(chan time.Time, 1)}
	timer.tfMap = make(map[*timerFun]bool)
	return timer
}

func (timer *Timer) AddFun(d time.Duration, fun func()) {
	if d <= 0 {
		fun()
	} else {
		rd := time.Now().Add(d).UnixNano()
		if timer.tf == nil {
			timer.tf = &timerFun{d: rd, fun: fun}
			timer.T = time.NewTimer(d)
		} else {
			if timer.tf.d > rd {
				timer.T.Stop()
				timer.tfMap[timer.tf] = true
				timer.tf = &timerFun{d: rd, fun: fun}
				timer.T = time.NewTimer(d)
			} else {
				timer.tfMap[&timerFun{d: rd, fun: fun}] = true
			}
		}
	}
}

func (timer *Timer) Run() {
	timer.tf.fun()
	if len(timer.tfMap) > 0 {
		now := time.Now().UnixNano()
		ls := make([]*timerFun, 0, len(timer.tfMap))
		var minTF *timerFun
		for tf, _ := range timer.tfMap {
			if tf.d <= now {
				ls = append(ls, tf)
			} else if minTF == nil || tf.d < minTF.d {
				minTF = tf
			}
		}
		if minTF != nil {
			delete(timer.tfMap, minTF)
			timer.tf = minTF
			timer.T = time.NewTimer(time.Duration(timer.tf.d - now))
		}
		for _, tf := range ls {
			delete(timer.tfMap, tf)
			tf.fun()
		}
	}
}

func (timer *Timer) Reset() {
	if timer.tf != nil {
		t := timer.T
		timer.T = &time.Timer{C: make(chan time.Time, 1)}
		timer.tf = nil
		timer.tfMap = make(map[*timerFun]bool)
		t.Stop()
	}
}

func (timer *Timer) Stop() {
	if timer.tf == nil {
		timer.T.Stop()
	}
}
