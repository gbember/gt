// mmap_test.go
package mmap

import (
	"testing"
)

func TestMMap(t *testing.T) {
	m := new(Map)
	m.gsize = 5
	m.id = 1
	m.maxVNum = 1000
	m.am = make(map[uint16]*area)
	for k := uint16(1); k <= 10; k++ {
		for i := uint16(1); i <= 10; i++ {
			a := new(area)
			a.id = i * k
			a.points = make([]point, 0, 4)
			a.allLines = make([]*line, 0, 4)
			a.lines = make(map[*line]bool)
			a.areaMap = make(map[*area]bool)
			p := point{x: float32(i) * 28, y: float32(k) * 28}
			a.points = append(a.points, p)

			p = point{x: float32(i+1) * 28, y: float32(k) * 28}
			a.points = append(a.points, p)
			l := &line{sp: a.points[0], ep: p}
			a.lines[l] = true
			a.allLines = append(a.allLines, l)

			p = point{x: float32(i+1) * 28, y: float32(k+1) * 28}
			a.points = append(a.points, p)
			l = &line{sp: a.points[1], ep: p}
			a.lines[l] = true
			a.allLines = append(a.allLines, l)

			p = point{x: float32(i) * 28, y: float32(k+1) * 28}
			a.points = append(a.points, p)
			l = &line{sp: a.points[2], ep: p}
			a.lines[l] = true
			a.allLines = append(a.allLines, l)

			l = &line{sp: a.points[4-1], ep: a.points[0]}
			a.lines[l] = true
			a.allLines = append(a.allLines, l)
			m.am[a.id] = a
		}
	}
	m.init()
	points := m.FindPath(point{x: 30, y: 30}, point{x: 250, y: 270})
	t.Log(points)
}
