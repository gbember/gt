package navmesh

//type pmap struct {
//	table   []point
//	buckets int
//	size    int
//}

//// round 到最近的2的倍数
//func minBuckets(v int) int {
//	v--
//	v |= v >> 1
//	v |= v >> 2
//	v |= v >> 4
//	v |= v >> 8
//	v |= v >> 16
//	v++
//	return v
//}

//func hashInt32(p point) int {
//	x = ((x >> 16) ^ x) * 0x45d9f3b
//	x = ((x >> 16) ^ x) * 0x45d9f3b
//	x = ((x >> 16) ^ x)
//	return x
//}

//func NewPMap(maxsize int) *pmap {
//	buckets := minBuckets(maxsize)
//	return &pmap{size: 0, buckets: buckets, table: make([]point, buckets)}
//}

//// TODO rehash策略
//func (m *pmap) Add(p point) {
//	num_probes, bucket_count_minus_one := 0, m.buckets-1
//	bucknum := hashInt32(int(d.Dealid)) & bucket_count_minus_one
//	for {
//		if m.table[bucknum].Dealid == 0 { // insert, 不支持放入ID为0的Deal
//			m.size += 1
//			m.table[bucknum] = d
//			return
//		}
//		if m.table[bucknum].Dealid == d.Dealid { // update
//			m.table[bucknum] = d
//			return
//		}
//		num_probes += 1 // Open addressing with Linear probing
//		bucknum = (bucknum + num_probes) & bucket_count_minus_one
//	}
//}

//func (m *pmap) HasPoint(p point) bool {
//	num_probes, bucket_count_minus_one := 0, m.buckets-1
//	bucknum := hashInt32(int(id)) & bucket_count_minus_one
//	for {
//		if m.table[bucknum].Dealid == id {
//			return m.table[bucknum], true
//		}
//		if m.table[bucknum].Dealid == 0 {
//			return m.table[bucknum], false
//		}
//		num_probes += 1
//		bucknum = (bucknum + num_probes) & bucket_count_minus_one
//	}
//}

//func (m *pmap) RM(p point) bool {
//	num_probes, bucket_count_minus_one := 0, m.buckets-1
//	bucknum := hashInt32(int(id)) & bucket_count_minus_one
//	for {
//		if m.table[bucknum].Dealid == id {
//			deal := m.table[bucknum]
//			num_probes += 1
//			bucknum1 := bucknum
//			bucknum2 := (bucknum + num_probes) & bucket_count_minus_one
//			for bucknum2 < m.buckets && m.table[bucknum2].Dealid != 0 {
//				m.table[bucknum1] = m.table[bucknum2]
//				bucknum1 = bucknum2
//				num_probes += 1
//				bucknum2 = (bucknum + num_probes) & bucket_count_minus_one
//			}
//			m.table[bucknum1].Dealid = 0
//			return deal, true
//		}
//		if m.table[bucknum].Dealid == 0 {
//			return m.table[bucknum], false
//		}
//		num_probes += 1
//		bucknum = (bucknum + num_probes) & bucket_count_minus_one
//	}
//}
