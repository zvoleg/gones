package registers

type ControllReg struct {
	value           byte
	baseNameTable   uint16
	incrementer     uint16
	spriteTable     uint16
	backgroundTable uint16
	spriteSize      uint16
	generateNmi     bool
}

func (r *ControllReg) Write(value byte) {
	r.value = value
	r.baseNameTable = r.getBaseNameTableAddress()
	r.incrementer = r.getVramIncrement()
	r.spriteTable = r.getSpriteTableAddress()
	r.backgroundTable = r.getBackgroundTableAddress()
	r.spriteSize = uint16(r.getSpriteSize())
	r.generateNmi = r.generateNmiOnVb()
}

func (r *ControllReg) getBaseNameTableAddress() uint16 {
	var address uint16
	switch r.value & 0x3 {
	case 0:
		address = 0x2000
	case 1:
		address = 0x2400
	case 2:
		address = 0x2800
	case 3:
		address = 0x2C00
	}
	return address
}

func (r *ControllReg) getVramIncrement() uint16 {
	var increment uint16
	switch (r.value >> 2) & 0x1 {
	case 0:
		increment = 1
	case 1:
		increment = 32
	}
	return increment
}

func (r *ControllReg) getSpriteTableAddress() uint16 {
	var address uint16
	switch (r.value >> 3) & 0x1 {
	case 0:
		address = 0x0000
	case 1:
		address = 0x1000
	}
	return address
}

func (r *ControllReg) getBackgroundTableAddress() uint16 {
	var address uint16
	switch (r.value >> 4) & 0x1 {
	case 0:
		address = 0x0000
	case 1:
		address = 0x1000
	}
	return address
}

func (r *ControllReg) getSpriteSize() int {
	var size int
	switch (r.value >> 5) & 0x1 {
	case 0:
		size = 8
	case 1:
		size = 16
	}
	return size
}

func (r *ControllReg) generateNmiOnVb() bool {
	return (r.value & 0x80) != 0
}

func (r *ControllReg) GenerateNmi() bool {
	return r.generateNmi
}

func (r *ControllReg) GetBackgroundTable() uint16 {
	return r.backgroundTable
}

func (r *ControllReg) Incrementer() uint16 {
	return r.incrementer
}

func (r *ControllReg) SpriteSize() uint16 {
	return r.spriteSize
}
