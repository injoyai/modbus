package modbus

import (
	"bytes"
	"encoding/binary"
	"math"
)

func ToBytes(any interface{}) ([]byte, error) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesBuffer, binary.BigEndian, any)
	return bytesBuffer.Bytes(), err
}

func BinToByte(bytes []bool) (result byte) {
	for _, b := range bytes {
		result *= 2
		if b {
			result++
		}
	}
	return
}

func ByteToBin(b byte) (list []bool) {
	list = make([]bool, 8)
	for i, _ := range list {
		x := uint8(math.Pow(2, float64(8-i-1)))
		if b >= x {
			list[8-i-1] = true
			b -= x
		}
	}
	return
}

func BytesToBin(bytes []byte) (list []bool) {
	for _, v := range bytes {
		list = append(list, ByteToBin(v)...)
	}
	return
}

// ReverseBool 倒序
func ReverseBool(bs []bool) []bool {
	x := make([]bool, len(bs))
	for i, v := range bs {
		x[len(bs)-i-1] = v
	}
	return x
}

// CoilsBytes 线圈转字节
func CoilsBytes(value []bool) (result []byte) {
	length := (len(value) + 7) / 8
	for i := 0; i < length; i++ {
		var x []bool
		if i == length-1 {
			x = ReverseBool(value[i*8:])
		} else {
			x = ReverseBool(value[i*8 : i*8+8])
		}
		result = append(result, BinToByte(x))
	}
	return append([]byte{byte(length)}, result...)
}

func CopyBytes(bytes []byte) []byte {
	result := []byte{}
	for _, v := range bytes {
		result = append(result, v)
	}
	return result
}
