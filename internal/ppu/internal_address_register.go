package ppu

type internalAddrReg struct {
	curValue uint16
	tmpValue uint16
	fineX    uint16
	latch    bool
}

func (reg *internalAddrReg) getFineY() uint16 {
	return reg.curValue >> 12
}

func (reg *internalAddrReg) getCoarseX() uint16 {
	return reg.curValue & 0x1F
}

func (reg *internalAddrReg) getCoarseY() uint16 {
	return (reg.curValue >> 5) & 0x1F
}

func (reg *internalAddrReg) incrementCoarseX() {
	if reg.curValue&0x001F == 31 {
		reg.curValue &= ^uint16(0x001F)
		reg.switchBit(0x0400)
	} else {
		reg.curValue += 1
	}
}

func (reg *internalAddrReg) incrementY() {
	if reg.curValue&0x7000 != 0x7000 {
		reg.curValue += 0x1000 // increment fine Y
	} else {
		reg.curValue &= ^uint16(0x7000) // clear fine Y
		coarseY := (reg.curValue & 0x03E0) >> 5
		switch coarseY {
		case 29:
			coarseY = 0
			reg.switchBit(0x0800)
		case 31:
			coarseY = 0
		default:
			coarseY += 1
		}
		reg.curValue &= ^uint16(0x03E00)
		reg.curValue |= coarseY << 5
	}
}

func (reg *internalAddrReg) copyHorizontalPosition() {
	mask := uint16(0x041F) // horizontal component of address
	reg.curValue &= ^mask  // clear bits
	reg.curValue |= reg.tmpValue & mask
}

func (reg *internalAddrReg) copyVerticalPosition() {
	mask := uint16(0x7BE0) // vertical component of address
	reg.curValue &= ^mask
	reg.curValue |= reg.tmpValue & mask
}

func (reg *internalAddrReg) setNameTable(data byte) {
	reg.tmpValue |= uint16((data & 0x3)) << 10
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
		reg.tmpValue &= ^(uint16(0x7) << 12) // clear bits before set
		fineY := data & 0x07
		reg.tmpValue |= fineY << 12
		reg.tmpValue &= ^(uint16(0x1F) << 5) // clear bits before set
		coarseY := data >> 3
		reg.tmpValue |= coarseY << 5
		reg.swapLatch()
	} else {
		reg.fineX = data & 0x07
		reg.tmpValue &= ^uint16(0x1F) // clear bits before set
		coarseX := data >> 3
		reg.tmpValue |= coarseX
		reg.swapLatch()
	}
}

func (reg *internalAddrReg) addressWrite(dataByte byte) {
	data := uint16(dataByte)
	if reg.latch {
		reg.tmpValue &= ^uint16(0xFF) // clear bits before set
		reg.tmpValue |= data
		reg.updateCurValue()
	} else {
		reg.tmpValue &= ^uint16(0x4000) // clear 14th bit
		reg.tmpValue &= ^uint16(0x3F00) // clear bits before set
		reg.tmpValue |= (data & 0x3F) << 8
	}
	reg.swapLatch()
}

func (reg *internalAddrReg) updateCurValue() {
	reg.curValue = reg.tmpValue
}

func (reg *internalAddrReg) increment(incrementer uint16) {
	reg.curValue += incrementer
}

func (reg *internalAddrReg) switchBit(bit uint16) {
	selectedBit := reg.curValue&bit != 0
	if selectedBit {
		reg.curValue &= ^uint16(bit) // clear bit
	} else {
		reg.curValue |= uint16(bit) // set bit
	}
}
