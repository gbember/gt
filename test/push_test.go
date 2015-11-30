// push_test.go
package test

import (
	"container/heap"
	"log"
	"math/rand"
	"sort"
	"testing"
	"time"
)

type intheap struct {
	is []int
}

func new_intheap(size int) *intheap {
	return &intheap{
		is: make([]int, 0, size),
	}
}

func (ih *intheap) Push(a int) {
	index := sort.Search(ih.Len()-1, func(i int) bool { return ih.is[i] < a })
	log.Println(index)
	ih.is = append(ih.is, a)
	if index != -1 {
		copy(ih.is[index+1:], ih.is[index:])
		log.Println(index, a)
		ih.is[index] = a
	}
	log.Println("===", ih)

	//	ih.is = append(ih.is, i)
	//	sindex := 0
	//	maxIndex := len(ih.is) - 2
	//	if maxIndex < 0 {
	//		return
	//	}
	//	eindex := maxIndex
	//	index := (eindex - sindex) / 2
	//	sort.Search()
	//	for ; sindex < eindex; index = (eindex - sindex) / 2 {
	//		time.Sleep(time.Millisecond * 10)
	//		log.Println(sindex, eindex, index)
	//		if i > ih.is[index] {
	//			eindex = index
	//		} else if i < ih.is[index] {
	//			sindex = index
	//		} else {
	//			break
	//		}

	//	}

	//	if i > ih.is[index] {
	//		copy(ih.is[index+1:], ih.is[index:])
	//		ih.is[index] = i
	//	} else {
	//		index++
	//		if index < maxIndex {
	//			copy(ih.is[index+1:], ih.is[index:])
	//			ih.is[index] = i
	//		}
	//	}
}

func (ih *intheap) Pop() int {
	length := len(ih.is) - 1
	a := ih.is[length]
	ih.is = ih.is[:length]
	return a
}

func (ih *intheap) Clear() {
	ih.is = ih.is[0:0]
}

func (ih *intheap) Len() int {
	return len(ih.is)
}

type intList []int

func (il intList) Len() int           { return len(il) }
func (il intList) Less(i, j int) bool { return il[i] < il[j] }
func (il intList) Swap(i, j int)      { il[i], il[j] = il[j], il[i] }
func (il *intList) Push(x interface{}) {
	*il = append(*il, x.(int))
}
func (ol *intList) Pop() interface{} {
	old := *ol
	length := len(old)
	x := old[length-1]
	*ol = old[:length-1]
	return x
}

var (
	num = 100000
	gil *intList
	gih *intheap
	ics []int
)

func init() {
	is := make([]int, 0, num)
	gil = (*intList)(&is)
	gih = new_intheap(num)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	max := r.Intn(10000) + 10000
	ics = make([]int, max, max)
	for i := 0; i < max; i++ {
		ics[i] = r.Int()
	}
}

func TestIntHeap(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.Ldate)
	ih := new_intheap(100)
	ih.Push(100)
	ih.Push(90)
	//	ih.Push(80)
	//	ih.Push(70)
	//	ih.Push(110)

	//	ih.Push(70)
	//	ih.Push(120)
	//	ih.Push(140)
	//	ih.Push(100)
	t.Log(ih)
}

func BenchmarkIntHeap(b *testing.B) {
	gih.Clear()
	for i := 0; i < b.N; i++ {
		gih.Push(ics[i%len(ics)])
	}
	b.Log(gih.Len())
}

func BenchmarkIntList(b *testing.B) {
	old := *gil
	*gil = old[0:0]
	for i := 0; i < b.N; i++ {
		heap.Push(gil, ics[i%len(ics)])
	}
	b.Log(gil.Len())
}
