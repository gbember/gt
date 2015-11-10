// area.go
package mmap

//区域
type area struct {
	id      uint16         //区域ID
	points  []point        //区域定点
	lines   []line         //区域不可穿过线段
	areaMap map[*area]bool //区域相邻区域
}

func loadArea(pk *Packet) (*area, error) {
	id, err := pk.readUint16()
	if err != nil {
		return nil, err
	}
	pointNum, err := pk.readUint16()
	if err != nil {
		return nil, err
	}
	a := new(area)
	a.id = id
	a.points = make([]point, 0, pointNum)
	a.lines = make([]line, 0, pointNum)
	a.areaMap = make(map[*area]bool)
	for ; pointNum > 0; pointNum-- {
		p := point{}
		p.x, err = pk.readFloat32()
		if err != nil {
			return nil, err
		}
		p.y, err = pk.readFloat32()
		if err != nil {
			return nil, err
		}
		a.points = append(a.points, p)
	}
	return a, nil
}


//是否穿过不能穿过的线
func (a *area) isCrossNoPassLine(l *line) bool {
	length := len(a.lines)
	for i := 0; i < length; i++ {
		if a.lines[i].isCross(l) {
			return false
		}
	}
	return true
}
