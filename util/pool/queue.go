// pool.go
package pool

type Queue struct {
	popPos  int
	pushPos int
	bs      []interface{}
}

func NewQueue(c int) *Queue {
	return &Queue{bs: make([]interface{}, c, c)}
}

func (this *Queue) Pop() (interface{}, bool) {
	if this.popPos == this.pushPos {
		return 0, false
	}
	v := this.bs[this.popPos]
	this.popPos++
	return v, true
}

func (this *Queue) Push(v interface{}) {
	this.bs[this.pushPos] = v
	this.pushPos++
	if this.pushPos >= len(this.bs) {
		this.pushPos = 0
	}
	if this.pushPos == this.popPos {
		l := 2 * len(this.bs)
		bs := make([]interface{}, l, l)
		copy(bs, this.bs)
		this.popPos = 0
		this.pushPos = len(this.bs)
		this.bs = bs
	}
}
