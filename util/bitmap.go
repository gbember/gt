// bitmap
//位标识功能
package util

import (
	"fmt"
)

type Bitmap struct {
	//保存实际的bit数据
	data []byte
	//指示该Bitmap的bit容量
	bitSize uint64
	//该Bitmap被设置为1的最大位置(方便遍历)
	maxPos uint64
}

func NewDefaultBitmap() *Bitmap {
	cap := uint64(1024)
	return NewBitmap(cap)
}

func NewBitmap(cap uint64) *Bitmap {
	bs := make([]byte, cap, cap)
	return &Bitmap{data: bs, bitSize: cap * 8}
}

func (this *Bitmap) MaxPos() uint64 {
	return this.maxPos
}

func (this *Bitmap) BitSize() uint64 {
	return this.bitSize
}

func (this *Bitmap) GetBit(offset uint64) (uint8, error) {
	index, pos := offset/8, offset%8

	if this.bitSize < offset {
		return 0, fmt.Errorf("offset(%d) out of bounds(%d)", offset, this.bitSize)
	}
	return (this.data[index] >> (pos - 1)) & 0x01, nil
}

func (this *Bitmap) SetBit(offset uint64, value uint8) bool {
	index, pos := offset/8, offset%8

	if this.bitSize < offset {
		return false
	}

	if value == 0 {
		// &^ 清位
		this.data[index] &^= 0x01 << pos
	} else {
		this.data[index] |= 0x01 << pos
		//设置记录为1的最大值
		if this.maxPos < offset {
			this.maxPos = offset
		}
	}
	return true
}
