package modbus

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func NewRTU(addr byte) Interface {
	return &client{
		slave: addr,
		model: RTU,
	}
}

func NewTCP(addr byte) Interface {
	return &client{
		slave: addr,
		model: TCP,
	}
}

type client struct {
	slave byte
	model string
}

// 1-bit access

func (this *client) ReadOutputCoils(address, quantity uint16) (results []byte, err error) {
	if quantity < 1 || quantity > 2000 {
		return nil, errors.New("超出范围(1-2000):" + strconv.Itoa(int(quantity)))
	}
	addressBytes, err := ToBytes(address)
	if err != nil {
		return nil, err
	}
	quantityBytes, err := ToBytes(quantity)
	if err != nil {
		return nil, err
	}
	data := append(addressBytes, quantityBytes...)
	return this.Encode(ReadCoils, data)
}

func (this *client) ReadInputCoils(address, quantity uint16) (results []byte, err error) {
	if quantity < 1 || quantity > 2000 {
		return nil, errors.New("超出范围(1-2000):" + strconv.Itoa(int(quantity)))
	}
	addressBytes, err := ToBytes(address)
	if err != nil {
		return nil, err
	}
	quantityBytes, err := ToBytes(quantity)
	if err != nil {
		return nil, err
	}
	data := append(addressBytes, quantityBytes...)
	return this.Encode(ReadDiscreteInputs, data)
}

func (this *client) WriteCoils(address uint16, value bool) (results []byte, err error) {
	addressBytes, err := ToBytes(address)
	if err != nil {
		return nil, err
	}
	dataBytes := []byte{0x00, 0x00}
	if value {
		dataBytes = []byte{0xFF, 0x00}
	}
	data := append(addressBytes, dataBytes...)
	return this.Encode(WriteCoils, data)
}

func (this *client) WriteMultipleCoils(address, quantity uint16, value []bool) (results []byte, err error) {
	if quantity > 1968 {
		return nil, errors.New("超出范围(0-1968):" + strconv.Itoa(int(quantity)))
	}
	addressBytes, err := ToBytes(address)
	if err != nil {
		return nil, err
	}
	quantityBytes, err := ToBytes(quantity)
	if err != nil {
		return nil, err
	}
	if int(quantity) != len(value) {
		return nil, errors.New("写入数据长度错误")
	}
	data := append(addressBytes, quantityBytes...)
	data = append(data, CoilsBytes(value)...)
	return this.Encode(WriteMultipleCoils, data)
}

// 16-bit access

func (this *client) ReadInputRegisters(address, quantity uint16) ([]byte, error) {
	addressBytes, err := ToBytes(address)
	if err != nil {
		return nil, err
	}
	quantityBytes, err := ToBytes(quantity)
	if err != nil {
		return nil, err
	}
	return this.Encode(ReadInputRegisters, append(addressBytes, quantityBytes...))
}

func (this *client) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	addressBytes, err := ToBytes(address)
	if err != nil {
		return nil, err
	}
	quantityBytes, err := ToBytes(quantity)
	if err != nil {
		return nil, err
	}
	return this.Encode(ReadHoldingRegisters, append(addressBytes, quantityBytes...))
}

func (this *client) WriteRegisters(address, value uint16) (results []byte, err error) {
	addressBytes, err := ToBytes(address)
	if err != nil {
		return nil, err
	}
	valueBytes, err := ToBytes(value)
	if err != nil {
		return nil, err
	}
	data := append(addressBytes, valueBytes...)
	return this.Encode(WriteRegisters, data)
}

func (this *client) WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error) {
	if len(value) != int(quantity)*2 {
		return nil, errors.New("写入数据长度错误")
	}
	if quantity > 120 || quantity < 1 {
		return nil, errors.New("超出寄存器范围(0x0001-0x0078)")
	}
	addressBytes, err := ToBytes(address)
	if err != nil {
		return nil, err
	}
	quantityBytes, err := ToBytes(quantity)
	if err != nil {
		return nil, err
	}
	data := append(addressBytes, quantityBytes...)
	data = append(data, 2*byte(quantity))
	data = append(data, value...)
	return this.Encode(WriteMultipleRegisters, data)
}

func (this *client) ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error) {
	return nil, errors.New("未实现")
}

func (this *client) MaskWriteRegisters(address, andMask, orMask uint16) (results []byte, err error) {
	return nil, errors.New("未实现")
}

func (this *client) ReadFIFOQueue(address uint16) (results []byte, err error) {
	return nil, errors.New("未实现")
}

// Encode													  -crc-
// RTU 04 00 01 00 04 >>> 					05 04 00 01 00 04 a1 8d
// TCP 04 00 01 00 04 >>> 01 00 00 00 00 06 05 04 00 01 00 04
func (this *client) Encode(control Control, data []byte) (result []byte, err error) {
	switch this.model {
	case RTU:
		result = EncodeRTU(this.slave, control, data)
	case TCP:
		result = EncodeTCP(this.slave, control, data)
	default:
		err = errors.New("未知Modbus类型:" + this.model)
	}
	return
}

// Decode 解析响应数据
func (this *client) Decode(bytes []byte) (frame Frame, err error) {
	if len(bytes) == 0 {
		return nil, errors.New("传感器未链接")
	}
	switch strings.ToUpper(this.model) {
	case RTU:
		frame, err = DecodeRTU(bytes)
	case TCP:
		frame, err = DecodeTCP(bytes)
	}
	if err != nil {
		return nil, err
	}
	if frame.GetSlave() != this.slave {
		return nil, fmt.Errorf("请求从站%v和响应从站%v不一致", this.slave, bytes[0])
	}
	return
}

// DecodeData 解析出数据域
// 050402011048ac >>> 0110
func (this *client) DecodeData(bytes []byte) ([]byte, error) {
	f, err := this.Decode(bytes)
	if err != nil {
		return nil, err
	}
	if int(f.GetData()[0]) != len(f.GetData()[1:]) {
		return nil, errors.New("数据域长度错误:" + f.HEX())
	}
	return f.GetData()[1:], nil
}
