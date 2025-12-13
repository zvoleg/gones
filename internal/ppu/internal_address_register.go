package ppu

type internalAddrReg struct {
	cur_value uint16
	tmp_value uint16
	fine_x    uint16
	latch     bool
}

func (reg *internalAddrReg) setNameTable(data byte) {
	reg.tmp_value |= uint16((data & 0x3)) << 10
}

func (reg *internalAddrReg) swapLatch() {
	reg.latch = !reg.latch
}

func (reg *internalAddrReg) resetLatch() {
	reg.latch = false
}

func (reg *internalAddrReg) scrollWrite(dataByte byte) {
	data := uint16(dataByte)
	if reg.latch {
		reg.tmp_value &= ^(uint16(0x7) << 12) // clear bits before set
		fine_y := data & 0x07
		reg.tmp_value |= fine_y << 12
		reg.tmp_value &= ^(uint16(0x1F) << 5) // clear bits before set
		coarse_y := data >> 3
		reg.tmp_value |= coarse_y << 5
		reg.swapLatch()
	} else {
		reg.fine_x = data & 0x07
		reg.tmp_value &= ^uint16(0x1F) // clear bits before set
		coarse_x := data >> 3
		reg.tmp_value |= coarse_x
		reg.swapLatch()
	}
}

func (reg *internalAddrReg) addressWrite(dataByte byte) {
	data := uint16(dataByte)
	if reg.latch {
		reg.tmp_value &= ^uint16(0xFF) // clear bits before set
		reg.tmp_value |= data
		reg.swapLatch()
		reg.updateCurValue()
	} else {
		reg.tmp_value &= ^(uint16(1) << 14)   // clear 14th bit
		reg.tmp_value &= ^(uint16(0x3F) << 8) // clear bits before set
		reg.tmp_value |= (data & 0x3F) << 8
		reg.swapLatch()
	}
}

func (reg *internalAddrReg) updateCurValue() {
	reg.cur_value = reg.tmp_value
}

func (reg *internalAddrReg) increment(incrementer uint16) {
	reg.cur_value += incrementer
}
