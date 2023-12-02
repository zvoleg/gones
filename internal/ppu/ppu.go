package ppu

import (
	"time"

	"github.com/zvoleg/gones/internal/cpu6502"
)

const freequencee float64 = 5369318.0
const clockTime float64 = 1.0 / freequencee
const clockTimeNs float64 = clockTime * 1000000000

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
}

type Ppu struct {
	patternTable [0x2000]byte
	nameTable    [0x1000]byte
	paletteRam   [0x0020]byte
	sram         [0x100]byte

	controllReg     *controllReg
	maskReg         *maskReg
	statusReg       *statusReg
	oamAddressReg   *oamAddressReg
	oamDataReg      *oamDataReg
	scrollRegister  *scrollRegister
	addressRegister *addressRegister

	bus           PpuBus
	cpuSignalLine chan cpu6502.Signal

	dmaEnabled   bool
	evenFrame    bool
	clockCounter int
}

func NewPpu(interruptLine chan cpu6502.Signal) Ppu {
	latch := false // scroll reg and address reg latch (selects LSB and HSB)
	controllReg := controllReg{}
	maskReg := maskReg{}
	statusReg := statusReg{latch: &latch}
	oamAddressReg := oamAddressReg{}
	oamDataReg := oamDataReg{}
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
	executionTimeNs := clockTimeNs * float64(time.Nanosecond)
	time.Sleep(time.Duration(executionTimeNs))
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
		ppu.cpuSignalLine <- cpu6502.Nmi
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
