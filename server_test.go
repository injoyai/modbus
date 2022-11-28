package modbus

import (
	"testing"
)

/*
01 00 00 00 00 06 01
*/
func TestNewServer(t *testing.T) {
	s := NewServer()
	s.Debug()
	if err := s.ListenTCP(503); err != nil {
		t.Error(err)
		return
	}
	testR := [2]byte{1, 2}
	s.SetHoldingRegisters(1, ReadWriteRegister{
		Read: func() ([2]byte, error) {
			return testR, nil
		},
		Write: func(bytes [2]byte) error {
			testR = bytes
			return nil
		},
	})
	s.SetInputRegisters(1, ReadWriteRegister{
		Read: func() ([2]byte, error) {
			return testR, nil
		},
		Write: func(bytes [2]byte) error {
			testR = bytes
			return nil
		},
	})
	testC := true
	s.SetCoils(2, ReadWriteCoils{
		Read: func() (bool, error) {
			return testC, nil
		},
		Write: func(b bool) error {
			testC = b
			return nil
		},
	})

	select {}
}
