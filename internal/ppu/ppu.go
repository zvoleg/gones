package ppu

import (
	"fmt"
)

type lineType int

const (
	Visible lineType = iota
	PostRender
	PreRender
)

type PpuBus interface {
	PpuRead(address uint16) byte
	PpuWrite(address uint16, data byte)
	ReadDmaByte(address uint16) byte
	GetMirroring() Mirroring
}

type Ppu struct {
	nameTable  [0x0800]byte
	paletteRam [0x0020]byte
	sram       [0x100]byte

	controllReg     *controllReg
	maskReg         *maskReg
	statusReg       *statusReg
	oamAddressReg   *oamAddressReg
	oamDataReg      *oamDataReg
	internalAddrReg *internalAddrReg
	// scrollRegister  *scrollRegister
	addressRegister *addressRegister

	dataBuffer byte

	bus PpuBus

	dmaEnabled     bool
	dmaByte        byte
	dmaClockWaiter bool

	interruptSignal bool

	evenFrame    bool
	clockCounter int
}

func NewPpu() Ppu {
	controllReg := controllReg{}
	maskReg := maskReg{}
	statusReg := statusReg{}
	oamAddressReg := oamAddressReg{}
	oamDataReg := oamDataReg{}
	internalAddrReg := internalAddrReg{}
	latch := false // scroll reg and address reg latch (selects LSB and HSB)
	// scrollReg := scrollRegister{latch: &latch}
	addressReg := addressRegister{latch: &latch}

	return Ppu{
		controllReg:     &controllReg,
		maskReg:         &maskReg,
		statusReg:       &statusReg,
		oamAddressReg:   &oamAddressReg,
		oamDataReg:      &oamDataReg,
		internalAddrReg: &internalAddrReg,
		// scrollRegister:  &scrollReg,
		addressRegister: &addressReg,
		bus:             nil,

		dmaEnabled:     false,
		dmaByte:        0,
		dmaClockWaiter: false,

		evenFrame:    false,
		clockCounter: 0,
	}
}

func (ppu *Ppu) DmaClockWaiter() bool {
	return ppu.dmaClockWaiter
}

func (ppu *Ppu) ResetDmaClockWaiter() {
	ppu.dmaClockWaiter = false
}

func (ppu *Ppu) InitBus(bus PpuBus) {
	ppu.bus = bus
}

func (ppu *Ppu) Clock() {
	if ppu.clockCounter == 0 && !ppu.evenFrame {
		ppu.clockCounter += 1
		return
	}

	lineNum := ppu.clockCounter / 341
	dotNum := ppu.clockCounter % 341
	line := getLineType(lineNum)

	if ppu.maskReg.renderingEnabled() {
		if line == Visible && dotNum < 256 {
			ppu.internalAddrReg.incrementCoarseX()
		}
		if line == Visible && dotNum == 256 {
			ppu.internalAddrReg.incrementY()
		}
		if dotNum == 257 { // copy horizontal position from tmp register
			ppu.internalAddrReg.copyHorizontalPosition()
		}
		if line == PreRender && (dotNum >= 280 && dotNum < 304) {
			ppu.internalAddrReg.copyVerticalPosition()
		}

		if line == Visible {
			//TODO fetch background pixel
			//TODO fetch plane pixel
			//TODO fetch sprite pixel
		}
	} else {
		// TODO make noise image?
	}

	if line == Visible || line == PreRender && dotNum >= 257 && dotNum <= 320 {
		ppu.oamAddressReg.value = 0
	}
	if lineNum == 241 && dotNum == 1 {
		ppu.statusReg.setStatusFlag(V, true)
		if ppu.controllReg.generateNmi {
			ppu.interruptSignal = true
		}
	}
	if line == PreRender && dotNum == 1 {
		ppu.statusReg.setStatusFlag(V, false)
		ppu.statusReg.setStatusFlag(S, false)
		ppu.statusReg.setStatusFlag(O, false)
		ppu.internalAddrReg.updateCurValue()
	}
	if ppu.clockCounter > 89342 {
		ppu.clockCounter = 0
		ppu.evenFrame = !ppu.evenFrame
		return
	}
	ppu.clockCounter += 1
}

func (ppu *Ppu) InterruptSignal() bool {
	signal := ppu.interruptSignal
	ppu.interruptSignal = false
	return signal
}

func getLineType(lineNum int) lineType {
	var line lineType
	if lineNum >= 0 && lineNum < 240 {
		line = Visible
	} else if lineNum >= 240 && lineNum < 261 {
		line = PostRender
	} else {
		line = PreRender
	}
	return line
}

func (ppu *Ppu) readRam(address uint16) byte {
	var data byte
	switch true {
	case address <= 0x1FFF:
		data = ppu.bus.PpuRead(address)
	case address >= 0x2000 && address <= 0x3EFF:
		switch ppu.bus.GetMirroring() {
		case Vertical:
			address = address & 0x07FF
			data = ppu.nameTable[address]
		case Horizontal:
			if address >= 0x2000 && address < 0x2400 {
				data = ppu.nameTable[address&0x03FF]
			}
			if address >= 0x2800 && address < 0x2C00 {
				data = ppu.nameTable[0x0400+(address&0x03FF)]
			}
			if address >= 0x2400 && address < 0x2800 {
				data = ppu.nameTable[address&0x03FF]
			}
			if address >= 0x2C00 && address < 0x3000 {
				data = ppu.nameTable[0x0400+(address&0x03FF)]
			}
		}
	case address >= 0x3F00 && address <= 0x3FFF:
		address = address & 0x1F
		switch address {
		case 0x04, 0x08, 0x0C:
			address = 0
		case 0x10, 0x14, 0x18, 0x1C:
			address = address & 0xF
		}
		data = ppu.paletteRam[address]
	default:
		fmt.Printf("Wrong address for writing into vram: 0x%04X\n", address)
	}
	return data
}

func (ppu *Ppu) writeVram(address uint16, data byte) {
	// fmt.Printf("Write into VRAM, address: %04X\n", address)
	switch true {
	case address <= 0x1FFF:
		ppu.bus.PpuWrite(address, data)
	case address >= 0x2000 && address <= 0x3EFF:
		switch ppu.bus.GetMirroring() {
		case Vertical:
			address = address & 0x07FF
			ppu.nameTable[address] = data
		case Horizontal:
			if address >= 0x2000 && address < 0x2400 {
				ppu.nameTable[address&0x03FF] = data
			}
			if address >= 0x2800 && address < 0x2C00 {
				ppu.nameTable[0x0400+(address&0x03FF)] = data
			}
			if address >= 0x2400 && address < 0x2800 {
				ppu.nameTable[address&0x03FF] = data
			}
			if address >= 0x2C00 && address < 0x3000 {
				ppu.nameTable[0x0400+(address&0x03FF)] = data
			}
		}
	case address >= 0x3F00 && address <= 0x3FFF:
		address = address & 0x1F
		switch address {
		case 0x10, 0x14, 0x18, 0x1C:
			mirrorAddress := address & 0xF
			ppu.paletteRam[address] = data
			ppu.paletteRam[mirrorAddress] = data
		default:
			ppu.paletteRam[address] = data
		}
	default:
		fmt.Printf("Wrong address for writing into vram: 0x%04X\n", address)
	}
}
