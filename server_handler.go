package modbus

import (
	"bufio"
	"io"
	"log"
)

type Handler func(f Frame) ([]byte, Control)

func (this *Server) handle(f Frame, w io.Writer) (err error) {
	origin := f.Copy()
	defer func() {
		if this.printHandler != nil {
			this.printHandler(origin, f)
		}
	}()
	switch f.GetControl() {
	case 1, 2, 3, 4, 5, 6, 15, 16:
		handler := this.Handler[f.GetControl()]
		if handler != nil {
			result, control := handler(f)
			if control != Success {
				f.SetControl(control)
			}
			f.SetData(result)
		}
	default:
		f.SetControl(IllegalFunction)
	}
	_, err = w.Write(f.Bytes())
	return err
}

// handler1 读线圈
func (this *Server) handler1(f Frame) ([]byte, Control) {
	return this.handlerReadCoils(f, this.Coils.Get)
}

// handler2 读离散输入
func (this *Server) handler2(f Frame) ([]byte, Control) {
	return this.handlerReadCoils(f, this.DiscreteInputs.Get)
}

// handler3 读保持寄存器
func (this *Server) handler3(f Frame) (result []byte, code Control) {
	return this.handlerReadRegister(f, this.HoldingRegisters.Get)
}

// handler4 读输入寄存器
func (this *Server) handler4(f Frame) (result []byte, code Control) {
	return this.handlerReadRegister(f, this.InputRegisters.Get)
}

// handler5 写一个线圈
// [53 45 0 0 0 6 1 5 0 2 255 0] //开
// [53 49 0 0 0 6 1 5 0 2 0 0]   //关
func (this *Server) handler5(f Frame) ([]byte, Control) {
	data := f.GetData()
	start := 256*uint16(data[0]) + uint16(data[1])
	value := data[2] == 255
	if start == 0 {
		return nil, IllegalAddress
	}
	err := this.Coils.Get(start, 1)[0].Write(value)
	if err != nil {
		return nil, DeviceFault
	}
	return data, Success
}

// handler6 写一个保持寄存器
func (this *Server) handler6(f Frame) ([]byte, Control) {
	data := f.GetData()
	start := 256*uint16(data[0]) + uint16(data[1])
	value := [2]byte{data[2], data[3]}
	if start == 0 {
		return nil, IllegalAddress
	}
	err := this.HoldingRegisters.Get(start, 1)[0].Write(value)
	if err != nil {
		return nil, DeviceFault
	}
	return data, Success
}

// handler15 写多个线圈
// [0 2 0 1 1 1]
func (this *Server) handler15(f Frame) ([]byte, Control) {
	data := f.GetData()
	if len(data) < 5 {
		return nil, IllegalAddress
	}
	start := 256*uint16(data[0]) + uint16(data[1])
	count := 256*uint16(data[2]) + uint16(data[3])
	list := BytesToBin(data[5:])
	if int(count*8) != len(list) || start == 0 {
		return nil, IllegalAddress
	}
	for i, v := range this.Coils.Get(start, count) {
		if err := v.Write(list[i]); err != nil {
			return nil, DeviceFault
		}
		data = data[2:]
	}
	return f.GetData(), Success
}

// handler16 写多个保持寄存器
// [0 1 0 1 2 1 3]
func (this *Server) handler16(f Frame) ([]byte, Control) {
	data := f.GetData()
	if len(data) < 5 {
		return nil, IllegalAddress
	}
	start := 256*uint16(data[0]) + uint16(data[1])
	count := 256*uint16(data[2]) + uint16(data[3])
	if count*2 != uint16(data[4]) || len(data[5:]) != int(data[4]) || start == 0 {
		return nil, IllegalAddress
	}
	data = data[5:]
	for _, v := range this.HoldingRegisters.Get(start, count) {
		if err := v.Write([2]byte{data[0], data[1]}); err != nil {
			return nil, DeviceFault
		}
		data = data[2:]
	}
	return f.GetData(), Success
}

// handlerReadCoils 读线圈
func (this *Server) handlerReadCoils(f Frame, fn func(uint16, uint16) []ReadWriteCoils) (result []byte, code Control) {
	data := f.GetData()
	start := 256*uint16(data[0]) + uint16(data[1])
	count := 256*uint16(data[2]) + uint16(data[3])
	value := []bool{}
	for _, v := range fn(start, count) {
		b, err := v.Read()
		if err != nil {
			return nil, IllegalFunction
		}
		value = append(value, b)
	}
	return CoilsBytes(value), Success
}

// handlerReadRegister 读寄存器
func (this *Server) handlerReadRegister(f Frame, fn func(uint16, uint16) []ReadWriteRegister) (result []byte, code Control) {
	data := f.GetData()
	start := 256*uint16(data[0]) + uint16(data[1])
	count := 256*uint16(data[2]) + uint16(data[3])
	for _, v := range fn(start, count) {
		bs, err := v.Read()
		if err != nil {
			return nil, DeviceFault
		}
		result = append(result, bs[0], bs[1])
	}
	result = append([]byte{byte(len(result))}, result...)
	return
}

func (this *Server) defaultPrintHandler(origin, result Frame) {
	if this.debug {
		log.Printf("[Modbus][%s] %s >>> %s ", origin.Type(), origin.HEX(), result.HEX())
	}
}

// ReadWithRTU 根据RTU数据格式读取数据,按字节(传输)读取
func ReadWithRTU(buf *bufio.Reader) (*RTUFrame, error) {
	bytes := []byte(nil)
	for {
		b, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, b)
		if len(bytes) >= 8 {
			f, err := DecodeRTU(bytes)
			if err == nil {
				return f, nil
			}
		}
	}
}

// ReadWithTCP 根据TCP数据格式读取数据
func ReadWithTCP(reader io.Reader) (*TCPFrame, error) {
	buff := make([]byte, 512)
	length, err := reader.Read(buff)
	if err != nil {
		return nil, err
	}
	buff = buff[:length]
	return DecodeTCP(buff)
}
