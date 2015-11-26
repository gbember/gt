// math.go
package navmesh

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

func maxmin(a, b int64) (int64, int64) {
	if a > b {
		return a, b
	}
	return b, a
}
