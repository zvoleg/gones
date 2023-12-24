package ppu

import (
	"fmt"

	"github.com/zvoleg/gones/internal/cpu6502"
)

// const freequencee float64 = 5369318.0
// const clockTime float64 = 1.0 / freequencee

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
	scrollRegister  *scrollRegister
	addressRegister *addressRegister

	dataBuffer byte

	bus           PpuBus
	cpuSignalLine chan cpu6502.Signal

	dmaEnabled   bool
	evenFrame    bool
	clockCounter int
}

func NewPpu(interruptLine chan cpu6502.Signal) Ppu {
	controllReg := controllReg{}
	maskReg := maskReg{}
	statusReg := statusReg{}
	oamAddressReg := oamAddressReg{}
	oamDataReg := oamDataReg{}
	latch := false // scroll reg and address reg latch (selects LSB and HSB)
	scrollReg := scrollRegister{latch: &latch}
	addressReg := addressRegister{latch: &latch}

	return Ppu{
		controllReg:     &controllReg,
		maskReg:         &maskReg,
		statusReg:       &statusReg,
		oamAddressReg:   &oamAddressReg,
		oamDataReg:      &oamDataReg,
		scrollRegister:  &scrollReg,
		addressRegister: &addressReg,
		bus:             nil,
		cpuSignalLine:   interruptLine,
		evenFrame:       false,
		clockCounter:    0,
	}
}

func (ppu *Ppu) InitBus(bus PpuBus) {
	ppu.bus = bus
}

func (ppu *Ppu) Clock() {
	if ppu.dmaEnabled {
		ppu.DmaClock()
	}
	if ppu.clockCounter == 0 && !ppu.evenFrame {
		ppu.clockCounter += 1
		return
	}

	lineNum := ppu.clockCounter / 341
	dotNum := ppu.clockCounter % 341
	line := getLineType(lineNum)

	if line == Visible || line == PreRender && dotNum >= 257 && dotNum <= 320 {
		ppu.oamAddressReg.value = 0
	}
	if lineNum == 241 && dotNum == 1 {
		ppu.statusReg.setStatusFlag(V, true)
		if ppu.controllReg.generateNmi {
			ppu.cpuSignalLine <- cpu6502.Nmi
		}
	}
	if lineNum == 261 && dotNum == 1 {
		ppu.statusReg.setStatusFlag(V, false)
		ppu.statusReg.setStatusFlag(S, false)
		ppu.statusReg.setStatusFlag(O, false)
	}
	if ppu.clockCounter > 89342 {
		ppu.clockCounter = 0
		ppu.evenFrame = !ppu.evenFrame
		return
	}
	ppu.clockCounter += 1
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

func (ppu *Ppu) readVram() byte {
	address := ppu.addressRegister.value
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
		fmt.Printf("Wrong address for writing into vram: %04X\n", address)
	}
	return data
}

func (ppu *Ppu) writeVram(data byte) {
	address := ppu.addressRegister.value
	fmt.Printf("Write into VRAM, address: %04X\n", address)
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
		fmt.Printf("Wrong address for writing into vram: %04X\n", address)
	}
}
