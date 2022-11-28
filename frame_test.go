package modbus

import (
	"encoding/hex"
	"testing"
)

func TestNewTCPFrame(t *testing.T) {
	{
		f, err := DecodeTCP([]byte{1, 0, 0, 0, 0, 6, 1, 3, 0, 1, 0, 10})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(f)
	}
	{
		f, err := DecodeTCP([]byte{1, 0, 0, 0, 0, 6, 1, 3})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(f)
	}
}

func TestNewRTUFrame(t *testing.T) {
	bs, err := hex.DecodeString("01030400000063ba1a")
	if err != nil {
		t.Error(err)
		return
	}
	f, err := DecodeRTU(bs)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(hex.EncodeToString(f.GetData()))
}

func TestNewTCPFrame1(t *testing.T) {
	bs, err := hex.DecodeString("01030400000063ba1a")
	if err != nil {
		t.Error(err)
		return
	}
	f, err := DecodeTCP(bs)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(hex.EncodeToString(f.GetData()))
}
