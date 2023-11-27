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
}

type Ppu struct {
	ControllReg   *ControllReg
	MaskReg       *MaskReg
	StatusReg     *StatusReg
	OamAddressReg *OamAddressReg
	OamDataReg    *OamDataReg

	bus           PpuBus
	interruptLine chan cpu6502.Interrupt

	evenFrame    bool
	clockCounter int
}

func NewPpu(interruptLine chan cpu6502.Interrupt) Ppu {
	controllReg := ControllReg{}
	maskReg := MaskReg{}
	statusReg := StatusReg{}
	oamAddressReg := OamAddressReg{}
	oamDataReg := OamDataReg{}

	return Ppu{
		ControllReg:   &controllReg,
		MaskReg:       &maskReg,
		StatusReg:     &statusReg,
		OamAddressReg: &oamAddressReg,
		OamDataReg:    &oamDataReg,
		bus:           nil,
		interruptLine: interruptLine,
		evenFrame:     false,
		clockCounter:  0,
	}
}

func (ppu *Ppu) InitBus(bus PpuBus) {
	ppu.bus = bus
}

func (ppu *Ppu) Clock() {
	executionTimeNs := clockTimeNs * float64(time.Nanosecond)
	time.Sleep(time.Duration(executionTimeNs))
	if ppu.clockCounter == 0 && !ppu.evenFrame {
		ppu.clockCounter += 1
		return
	}

	lineNum := ppu.clockCounter / 341
	dotNum := ppu.clockCounter % 341
	line := getLineType(lineNum)

	if line == Visible || line == PreRender && dotNum >= 257 && dotNum <= 320 {
		ppu.OamAddressReg.value = 0
	}
	if lineNum == 241 && dotNum == 1 {
		ppu.StatusReg.setStatusFlag(V, true)
		ppu.interruptLine <- cpu6502.Nmi
	}
	if lineNum == 261 && dotNum == 1 {
		ppu.StatusReg.setStatusFlag(V, false)
		ppu.StatusReg.setStatusFlag(S, false)
		ppu.StatusReg.setStatusFlag(O, false)
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
