package registers

const (
	maskCoarseX    uint16 = 0x001F
	maskCoarseY    uint16 = 0x03E0
	maskFineY      uint16 = 0x7000
	maskHorizontal uint16 = 0x041F
	maskVertical   uint16 = 0x7BE0
)

type InternalAddrReg struct {
	curValue uint16
	tmpValue uint16
	fineX    uint16
	latch    bool
}

func (reg *InternalAddrReg) GetAddress() uint16 {
	return reg.curValue
}

func (reg *InternalAddrReg) GetFineY() uint16 {
	return (reg.curValue >> 12) & 0x7
}

func (reg *InternalAddrReg) GetFineX() uint16 {
	return reg.fineX
}

func (reg *InternalAddrReg) GetCoarseX() uint16 {
	return reg.curValue & 0x1F
}

func (reg *InternalAddrReg) GetCoarseY() uint16 {
	return (reg.curValue >> 5) & 0x1F
}

func (reg *InternalAddrReg) IncrementCoarseX() {
	if reg.curValue&maskCoarseX == 31 {
		reg.curValue &= ^maskCoarseX
		reg.curValue ^= 0x0400
	} else {
		reg.curValue += 1
	}
}

func (reg *InternalAddrReg) IncrementY() {
	if reg.curValue&maskFineY != maskFineY {
		reg.curValue += 0x1000 // increment fine Y
	} else {
		reg.curValue &= ^maskFineY // clear fine Y
		coarseY := (reg.curValue & maskCoarseY) >> 5
		switch coarseY {
		case 29:
			coarseY = 0
			reg.curValue ^= 0x0800
		case 31:
			coarseY = 0
		default:
			coarseY += 1
		}
		reg.curValue &= ^maskCoarseY
		reg.curValue |= coarseY << 5
	}
}

func (reg *InternalAddrReg) CopyHorizontalPosition() {
	reg.curValue &= ^maskHorizontal // clear bits
	reg.curValue |= reg.tmpValue & maskHorizontal
}

func (reg *InternalAddrReg) CopyVerticalPosition() {
	reg.curValue &= ^maskVertical
	reg.curValue |= reg.tmpValue & maskVertical
}

func (reg *InternalAddrReg) SetNameTable(data byte) {
	reg.tmpValue &= ^(uint16(0x3) << 10)
	reg.tmpValue |= uint16(data&0x3) << 10
}

func (reg *InternalAddrReg) ResetLatch() {
	reg.latch = false
}

func (reg *InternalAddrReg) ScrollWrite(dataByte byte) {
	data := uint16(dataByte)
	if reg.latch {
		reg.tmpValue &= ^maskFineY // clear bits before set
		fineY := data & 0x07
		reg.tmpValue |= fineY << 12
		reg.tmpValue &= ^maskCoarseY // clear bits before set
		coarseY := data >> 3
		reg.tmpValue |= coarseY << 5
	} else {
		reg.fineX = data & 0x07
		reg.tmpValue &= ^maskCoarseX // clear bits before set
		coarseX := data >> 3
		reg.tmpValue |= coarseX
	}
	reg.swapLatch()
}

func (reg *InternalAddrReg) AddressWrite(dataByte byte) {
	data := uint16(dataByte)
	if reg.latch {
		reg.tmpValue &= uint16(0xFF00) // clear bits before set
		reg.tmpValue |= data
		reg.UpdateCurValue()
	} else {
		reg.tmpValue &= uint16(0x80FF) // clear bits before set
		reg.tmpValue |= (data & 0x3F) << 8
	}
	reg.swapLatch()
}

func (reg *InternalAddrReg) UpdateCurValue() {
	reg.curValue = reg.tmpValue
}

func (reg *InternalAddrReg) Increment(incrementer uint16) {
	reg.curValue += incrementer
}

func (reg *InternalAddrReg) swapLatch() {
	reg.latch = !reg.latch
}
