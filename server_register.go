package modbus

// ReadWriteRegister 寄存器读写
type ReadWriteRegister struct {
	Read  func() ([2]byte, error)
	Write func([2]byte) error
}

type Register [65535]ReadWriteRegister

func (this Register) Get(start, count uint16) (result []ReadWriteRegister) {
	for i := start; i < start+count; i++ {
		if i < 65535 {
			l := this[i]
			if l.Read == nil {
				l.Read = func() (result [2]byte, err error) { return }
			}
			if l.Write == nil {
				l.Write = func([2]byte) (err error) { return }
			}
			result = append(result, l)
		} else {
			result = append(result, ReadWriteRegister{
				Read:  func() (result [2]byte, err error) { return },
				Write: func([2]byte) (err error) { return },
			})
		}
	}
	return
}

// ReadWriteCoils 线圈读写
type ReadWriteCoils struct {
	Read  func() (bool, error)
	Write func(bool) error
}

type Coils [65535]ReadWriteCoils

func (this Coils) Get(start, count uint16) (result []ReadWriteCoils) {
	for i := start; i < start+count; i++ {
		if i < 65535 {
			l := this[i]
			if l.Read == nil {
				l.Read = func() (result bool, err error) { return }
			}
			if l.Write == nil {
				l.Write = func(bool) (err error) { return }
			}
			result = append(result, l)
		} else {
			result = append(result, ReadWriteCoils{
				Read:  func() (result bool, err error) { return },
				Write: func(bool) (err error) { return },
			})
		}
	}
	return
}
