package modbus

type Interface interface {
	// 1-bit access

	// ReadOutputCoils reads from 1 to 2000 contiguous status of coils in a
	// remote device and returns coil status.
	ReadOutputCoils(address, quantity uint16) (results []byte, err error)
	// ReadInputCoils reads from 1 to 2000 contiguous status of
	// discrete inputs in a remote device and returns input status.
	ReadInputCoils(address, quantity uint16) (results []byte, err error)
	// WriteCoils write a single output to either ON or OFF in a
	// remote device and returns output value.
	WriteCoils(address uint16, value bool) (results []byte, err error)
	// WriteMultipleCoils forces each coil in a sequence of coils to either
	// ON or OFF in a remote device and returns quantity of outputs.
	WriteMultipleCoils(address, quantity uint16, value []bool) (results []byte, err error)

	// 16-bit access

	// ReadInputRegisters
	// 请求 PDU
	// |功能码 		|1个字节 		|0x10
	// |起始地址 	|2个字节 		|0x0000 至 0xFFFF
	// |寄存器数量 	|2个字节 		|0x0001 至 0x0078
	// |字节数 		|1个字节 		|2×N*
	// |寄存器值 	|N*×2个字节 		|值
	// *N＝寄存器数量
	// 响应 PDU
	// |功能码   	|1个字节 		|0x10
	// |起始地址    	|2个字节 		|0x0000 至 0xFFFF
	// |寄存器数量  	|2个字节 		|1 至 123（0x7B）
	// 错误
	// |差错码 		|1个字节 		|0x90
	// |异常码 		|1个字节 		|01 或 02 或 03 或 04
	// ReadInputRegisters reads from 1 to 125 contiguous input registers in
	// a remote device and returns input registers.
	ReadInputRegisters(address, quantity uint16) (results []byte, err error)
	// ReadHoldingRegisters reads the contents of a contiguous block of
	// holding registers in a remote device and returns register value.
	ReadHoldingRegisters(address, quantity uint16) (results []byte, err error)
	// WriteRegisters writes a single holding register in a remote
	// device and returns register value.
	WriteRegisters(address, value uint16) (results []byte, err error)
	// WriteMultipleRegisters writes a block of contiguous registers
	// (1 to 123 registers) in a remote device and returns quantity of
	// registers.
	WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error)
	// ReadWriteMultipleRegisters performs a combination of one read
	// operation and one write operation. It returns read registers value.
	ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error)
	// MaskWriteRegisters modify the contents of a specified holding
	// register using a combination of an AND mask, an OR mask, and the
	// register's current contents. The function returns
	// AND-mask and OR-mask.
	MaskWriteRegisters(address, andMask, orMask uint16) (results []byte, err error)
	//ReadFIFOQueue reads the contents of a First-In-First-Out (FIFO) queue
	// of register in a remote device and returns FIFO value register.
	//ReadFIFOQueue(address uint16) (results []byte, err error)

	Encode(Control, []byte) ([]byte, error)
	Decode([]byte) (Frame, error)
	DecodeData([]byte) ([]byte, error)
}
