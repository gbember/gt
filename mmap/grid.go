// grid.go
package mmap

//得到点所在的格子编号
//p 点
//gsize 正方形格子的宽
//maxVNum 横向最大格子数
func getGridNum(p point, gsize int32, maxVNum int32) int32 {
	return p.x/gsize + 1 + p.y/gsize*maxVNum
}
