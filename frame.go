package modbus

import (
	"encoding/hex"
	"errors"
	"fmt"
)

type Frame interface {

	// Type 类型2种 "TCP" or "RTU"
	Type() string

	// Copy 复制
	Copy() Frame

	// Bytes 根据当前数据,(实时)编码成字节
	Bytes() []byte

	// GetSlave 获取从站地址
	GetSlave() byte

	// SetSlave 设置从站地址
	SetSlave(slave byte)

	// GetControl 获取控制码
	GetControl() Control

	// SetControl 设置控制码
	SetControl(control Control)

	// GetData 获取数据域
	GetData() []byte

	// SetData 设置数据域
	SetData(data []byte)

	// HEX 十六进制字符
	HEX() string
}

type RTUFrame struct {
	Slave   byte
	Control Control
	Data    []byte
	CRC     []byte
}

func (this *RTUFrame) Type() string {
	return RTU
}

func (this *RTUFrame) Copy() Frame {
	x := *this
	return &x
}

func (this *RTUFrame) HEX() string {
	return hex.EncodeToString(this.Bytes())
}

func (this *RTUFrame) Bytes() []byte {
	bytes := []byte(nil)
	bytes = append(bytes, this.Slave, this.Control.Byte())
	bytes = append(bytes, this.Data...)
	bytes = append(bytes, CRC(bytes)...)
	return bytes
}

func (this *RTUFrame) GetSlave() byte {
	return this.Slave
}

func (this *RTUFrame) SetSlave(slave byte) {
	this.Slave = slave
}

func (this *RTUFrame) GetControl() Control {
	return this.Control
}

func (this *RTUFrame) SetControl(control Control) {
	this.Control = control
}

func (this *RTUFrame) GetData() []byte {
	return this.Data
}

func (this *RTUFrame) SetData(data []byte) {
	this.Data = data
}

func EncodeRTU(slave byte, control Control, data []byte) []byte {
	f := &RTUFrame{
		Slave:   slave,
		Control: control,
		Data:    data,
	}
	return f.Bytes()
}

func DecodeRTU(bytes []byte) (*RTUFrame, error) {
	length := len(bytes)
	if length < 5 {
		return nil, errors.New("数据长度异常(小于5):" + hex.EncodeToString(bytes))
	}
	crc := CRC(bytes[:length-2])
	if len(crc) != 2 || (crc[0] != bytes[length-2:][0] && crc[1] != bytes[length-2:][1]) {
		return nil, errors.New(fmt.Sprintf("crc校验错误:%s", hex.EncodeToString(bytes)))
	}
	if hex.EncodeToString(CRC(bytes[:len(bytes)-2])) != hex.EncodeToString(bytes[len(bytes)-2:]) {
		return nil, fmt.Errorf("CRC校验结果不一致%s", hex.EncodeToString(bytes))
	}
	bytes = bytes[:len(bytes)-2]
	return &RTUFrame{
		Slave:   bytes[0],
		Control: Control(bytes[1]),
		Data:    bytes[2:],
		CRC:     crc,
	}, nil
}

type TCPFrame struct {
	Order    [2]byte //序号,原路返回
	Protocol [2]byte //协议,原路返回
	Length   [2]byte //数据长度
	Slave    byte    //从站地址
	Control  Control //控制码
	Data     []byte  //数据域
}

func (this *TCPFrame) Type() string {
	return TCP
}

func (this *TCPFrame) Copy() Frame {
	x := *this
	return &x
}

func (this *TCPFrame) HEX() string {
	return hex.EncodeToString(this.Bytes())
}

func (this *TCPFrame) Bytes() []byte {
	bytes := []byte(nil)
	bytes = append(bytes, this.Order[0], this.Order[1])
	bytes = append(bytes, this.Protocol[0], this.Protocol[1])
	bytes = append(bytes, this.Length[0], this.Length[1])
	bytes = append(bytes, this.Slave, this.Control.Byte())
	bytes = append(bytes, this.Data...)
	return bytes
}

func (this *TCPFrame) GetSlave() byte {
	return this.Slave
}

func (this *TCPFrame) SetSlave(slave byte) {
	this.Slave = slave
}

func (this *TCPFrame) GetControl() Control {
	return this.Control
}

func (this *TCPFrame) SetControl(control Control) {
	this.Control = control
}

func (this *TCPFrame) GetData() []byte {
	return this.Data
}

func (this *TCPFrame) SetData(data []byte) {
	this.Data = data
	length := len(this.Data) + 2
	this.Length = [2]byte{byte(length / 256), byte(length % 256)}
}

func EncodeTCP(slave byte, control Control, data []byte) []byte {
	f := &TCPFrame{
		Order:    [2]byte{0x01, 0x00},
		Protocol: [2]byte{},
		Length:   [2]byte{0x00, 0x06},
		Slave:    slave,
		Control:  control,
		Data:     data,
	}
	return f.Bytes()
}

func DecodeTCP(bytes []byte) (*TCPFrame, error) {
	length := len(bytes)
	if length < 9 {
		return nil, errors.New("数据长度异常(小于9):" + hex.EncodeToString(bytes))
	}
	f := &TCPFrame{
		Order:    [2]byte{bytes[0], bytes[1]},
		Protocol: [2]byte{bytes[2], bytes[3]},
		Length:   [2]byte{bytes[4], bytes[5]},
		Slave:    bytes[6],
		Control:  Control(bytes[7]),
		Data:     bytes[8:],
	}
	if f.Control.Byte() > 0x80 {
		return f, f.Control
	}
	if int(f.Length[0])*256+int(f.Length[1]) != len(f.Data)+2 {
		return f, errors.New("数据长度错误:" + f.HEX())
	}
	return f, nil
}
