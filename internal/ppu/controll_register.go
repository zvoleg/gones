package ppu

type controllReg struct {
	value           byte
	baseNameTable   uint16
	incrementer     uint16
	spriteTable     uint16
	backgroundTable uint16
	spriteSize      uint16
	generateNmi     bool
}

func (r *controllReg) write(value byte) {
	r.value = value
	r.baseNameTable = r.getBaseNameTableAddress()
	r.incrementer = r.getVramIncrement()
	r.spriteTable = r.getSpriteTableAddress()
	r.backgroundTable = r.getBackgroundTableAddress()
	r.spriteSize = uint16(r.getSpriteSize())
	r.generateNmi = r.generateNmiOnVb()
}

func (r *controllReg) getBaseNameTableAddress() uint16 {
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

func (r *controllReg) getVramIncrement() uint16 {
	var increment uint16
	switch (r.value >> 2) & 0x1 {
	case 0:
		increment = 1
	case 1:
		increment = 32
	}
	return increment
}

func (r *controllReg) getSpriteTableAddress() uint16 {
	var address uint16
	switch (r.value >> 3) & 0x1 {
	case 0:
		address = 0x0000
	case 1:
		address = 0x1000
	}
	return address
}

func (r *controllReg) getBackgroundTableAddress() uint16 {
	var address uint16
	switch (r.value >> 4) & 0x1 {
	case 0:
		address = 0x0000
	case 1:
		address = 0x1000
	}
	return address
}

func (r *controllReg) getSpriteSize() int {
	var size int
	switch (r.value >> 5) & 0x1 {
	case 0:
		size = 8 * 8
	case 1:
		size = 8 * 16
	}
	return size
}

func (r *controllReg) generateNmiOnVb() bool {
	return (r.value & 0x80) != 0
}
