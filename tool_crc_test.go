package modbus

import (
	"testing"
)

func TestCRC(t *testing.T) {
	bs := []byte{0, 1, 2, 3}
	t.Log(CRC(bs))
}
