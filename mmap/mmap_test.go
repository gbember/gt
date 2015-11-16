// mmap_test.go
package mmap

import (
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

func TestMMap(t *testing.T) {
	m := loadMap_1()

	fn := "tt.cpuprof"
	f, err := os.Create(fn)
	isProf := false
	if err == nil {
		err = pprof.StartCPUProfile(f)
		if err == nil {
			isProf = true
		}
	}

	points, ok := m.FindPath(point{x: 101, y: 86}, point{x: 328, y: 324})
	if ok {
		t.Log(points)
		startTime := time.Now()
		var max int64 = 10000
		for i := max; i > 0; i-- {
			m.FindPath(point{x: 101, y: 86}, point{x: 328, y: 324})
		}
		td := time.Since(startTime)
		t.Log(td)
		t.Log(td.Nanoseconds() / max)
	} else {
		t.Log("终点不可达")
	}

	if isProf {
		pprof.StopCPUProfile()
	}
}

func addArea(am map[uint32]*area, id uint32, ps ...point) {
	a := new(area)
	a.id = id
	length := len(ps)
	if length < 3 {
		panic("area point num must lager 3")
	}
	a.points = ps
	a.allLines = make([]*line, 0, length)
	a.lineMap = make(map[*line]bool)
	a.areaMap = make(map[*line]*area)

	p := ps[0]
	l := &line{ps[length-1], p}
	a.allLines = append(a.allLines, l)

	for i := 1; i < length; i++ {
		l = &line{p, ps[i]}
		a.allLines = append(a.allLines, l)
		p = ps[i]
	}

	am[id] = a

}

func loadMap_1() *Map {
	m := new(Map)
	m.gsize = 5
	m.id = 1
	m.maxVNum = 1000
	m.am = make(map[uint32]*area)

	addArea(m.am, 1, point{75, 95}, point{118, 97}, point{188, 50}, point{190, 23})
	addArea(m.am, 2, point{190, 23}, point{188, 50}, point{294, 59}, point{307, 39}, point{207, 22})
	addArea(m.am, 3, point{270, 22}, point{307, 39}, point{318, 2}, point{282, 3})
	addArea(m.am, 4, point{307, 39}, point{294, 59}, point{354, 101}, point{381, 65})
	addArea(m.am, 5, point{381, 65}, point{354, 101}, point{408, 97}, point{364, 54}, point{566, 2}, point{476, 3})
	addArea(m.am, 6, point{408, 97}, point{354, 101}, point{380, 134}, point{425, 125})
	addArea(m.am, 7, point{425, 125}, point{380, 134}, point{368, 182}, point{420, 186})
	addArea(m.am, 8, point{420, 186}, point{368, 182}, point{349, 216}, point{403, 217})
	addArea(m.am, 9, point{403, 217}, point{349, 216}, point{291, 280}, point{347, 291})
	addArea(m.am, 10, point{347, 291}, point{291, 280}, point{246, 318}, point{247, 340})

	m.init()
	return m
}
