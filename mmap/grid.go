// grid.go
package mmap

//得到点所在的格子编号
//p 点
//gsize 正方形格子的宽
//maxVNum 横向最大格子数
func getGridNum(p point, gsize int, maxVNum int) int {
	return int(p.x)/gsize + 1 + (int(p.y)/gsize+1)*maxVNum
}

//得到线段所经过的格子编号列表
//l 线段
//gsize 正方形格子的宽
//maxVNum 横向最大格子数
func getGridNums(l *line, gsize int, maxVNum int) []int {
	gnum1 := getGridNum(l.sp, gsize, maxVNum)
	gnum2 := getGridNum(l.ep, gsize, maxVNum)
	if gnum1 == gnum2 {
		return []int{gnum1}
	}
	//TODO 未完成
	return []int{gnum1, gnum2}
}
