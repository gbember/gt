package test

// Deal的定义
type DealTiny struct {
	Dealid    int32
	Classid   int32
	Mttypeid  int32
	Bizacctid int32
	Isonline  bool
	Geocnt    int32
}

type DealMap struct {
	table   []DealTiny
	buckets int
	size    int
}

// round 到最近的2的倍数
func minBuckets(v int) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func hashInt32(x int) int {
	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = ((x >> 16) ^ x)
	return x
}

func NewDealMap(maxsize int) *DealMap {
	buckets := minBuckets(maxsize)
	return &DealMap{size: 0, buckets: buckets, table: make([]DealTiny, buckets)}
}

// TODO rehash策略
func (m *DealMap) Put(d DealTiny) {
	num_probes, bucket_count_minus_one := 0, m.buckets-1
	bucknum := hashInt32(int(d.Dealid)) & bucket_count_minus_one
	for {
		if m.table[bucknum].Dealid == 0 { // insert, 不支持放入ID为0的Deal
			m.size += 1
			m.table[bucknum] = d
			return
		}
		if m.table[bucknum].Dealid == d.Dealid { // update
			m.table[bucknum] = d
			return
		}
		num_probes += 1 // Open addressing with Linear probing
		bucknum = (bucknum + num_probes) & bucket_count_minus_one
	}
}

func (m *DealMap) Get(id int32) (DealTiny, bool) {
	num_probes, bucket_count_minus_one := 0, m.buckets-1
	bucknum := hashInt32(int(id)) & bucket_count_minus_one
	for {
		if m.table[bucknum].Dealid == id {
			return m.table[bucknum], true
		}
		if m.table[bucknum].Dealid == 0 {
			return m.table[bucknum], false
		}
		num_probes += 1
		bucknum = (bucknum + num_probes) & bucket_count_minus_one
	}
}

func (m *DealMap) Delete(id int32) (DealTiny, bool) {
	num_probes, bucket_count_minus_one := 0, m.buckets-1
	bucknum := hashInt32(int(id)) & bucket_count_minus_one
	for {
		if m.table[bucknum].Dealid == id {
			deal := m.table[bucknum]
			num_probes += 1
			bucknum1 := bucknum
			bucknum2 := (bucknum + num_probes) & bucket_count_minus_one
			for bucknum2 < m.buckets && m.table[bucknum2].Dealid != 0 {
				m.table[bucknum1] = m.table[bucknum2]
				bucknum1 = bucknum2
				num_probes += 1
				bucknum2 = (bucknum + num_probes) & bucket_count_minus_one
			}
			m.table[bucknum1].Dealid = 0
			return deal, true
		}
		if m.table[bucknum].Dealid == 0 {
			return m.table[bucknum], false
		}
		num_probes += 1
		bucknum = (bucknum + num_probes) & bucket_count_minus_one
	}
}
