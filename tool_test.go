package modbus

import (
	"encoding/hex"
	"testing"
)

func TestBinToByte(t *testing.T) {
	t.Log(BinToByte([]bool{true, false, true, true, true}))
}

func TestEntity(t *testing.T) {
	x := NewRTU(5)
	bytes, err := x.WriteMultipleCoils(19, 10, []bool{
		true, false, true, true, false, false, true, true, true, false})
	t.Error(err)
	t.Log(hex.EncodeToString(bytes))
}

func TestByteToBin(t *testing.T) {
	t.Log(ByteToBin(0xcd))
	t.Log(BytesToBin([]byte{0xcd, 0x01}))
}
