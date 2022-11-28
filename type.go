package modbus

const (

	// KeyOutputCoils 离散(线圈)输出（可读可写）
	// 对应PLC为：DO
	KeyOutputCoils = "0X"

	// KeyInputCoils 离散(线圈)输入（只读）
	// 对应PLC为：DI
	KeyInputCoils = "1X"

	// KeyInputRegisters 输入寄存器16位（只读）
	// 对应PLC为：AI
	KeyInputRegisters = "3X"

	// KeyHoldingRegisters 保持寄存器16位（可读可写）
	// 对应PLC为：AO
	KeyHoldingRegisters = "4X"

	// TCP ModbusTCP
	TCP = "TCP"

	// RTU ModbusRTU
	RTU = "RTU"
)
