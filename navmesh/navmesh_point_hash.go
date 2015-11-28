package navmesh

type pref_map struct {
	length int64
	flags  []bool
	maxX   int64
}

func hashcode(p Point, maxX int64) int64 {
	return hashInt32(p.X) + hashInt32(p.Y*maxX)
	//	return p.Y*(maxX+33) + p.X
}

func hashInt32(x int64) int64 {
	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = ((x >> 16) ^ x)
	return x
}

func makePrefMap(ps []Point, maxX int64) *pref_map {
	length := int64(len(ps))
	maxLength := length * 2
	for {
		if checkHakePrefMap(ps, maxX, length) {
			break
		}
		if length >= maxLength {
			return nil
		}
		length += 1
	}
	pm := new(pref_map)
	pm.length = length
	pm.flags = make([]bool, length, length)
	pm.maxX = maxX
	return pm
}

func checkHakePrefMap(ps []Point, maxX int64, length int64) bool {
	m := make(map[int64]bool)
	for i := 0; i < len(ps); i++ {
		hc := hashcode(ps[i], maxX) % length
		if m[hc] {
			return false
		}
		m[hc] = true
	}
	return true
}

func (pm *pref_map) isContainsPoint(p Point) bool {
	hc := hashcode(p, pm.maxX) % pm.length
	return pm.flags[hc]
}
func (pm *pref_map) addPoint(p Point) {
	hc := hashcode(p, pm.maxX) % pm.length
	pm.flags[hc] = true
}

func (pm *pref_map) clear() {
	length := pm.length
	for i := int64(0); i < length; i++ {
		pm.flags[i] = false
	}
}
