package ppu

import (
	"fmt"
	"math/rand"

	reg "github.com/zvoleg/gones/internal/ppu/internal/registers"
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

	screen     image // screen size 256x240, 4 byte per pixel (RGBA)
	screenPosX int
	screenPosY int

	planeDataBuffer planeDataBuffer
	shiftRegiseter  reg.ShiftRegiseter

	controllReg     *reg.ControllReg
	maskReg         *reg.MaskReg
	statusReg       *reg.StatusReg
	oamAddressReg   *reg.OamAddressReg
	oamDataReg      *reg.OamDataReg
	internalAddrReg *reg.InternalAddrReg

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
	controllReg := reg.ControllReg{}
	maskReg := reg.MaskReg{}
	statusReg := reg.StatusReg{}
	oamAddressReg := reg.OamAddressReg{}
	oamDataReg := reg.OamDataReg{}
	internalAddrReg := reg.InternalAddrReg{}

	return Ppu{
		screen: newImage(256, 240),

		controllReg:     &controllReg,
		maskReg:         &maskReg,
		statusReg:       &statusReg,
		oamAddressReg:   &oamAddressReg,
		oamDataReg:      &oamDataReg,
		internalAddrReg: &internalAddrReg,
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

	if ppu.maskReg.RenderingEnabled() {
		if line == Visible && dotNum < 256 {
			colorId := ppu.readRam(0x3F00) // fetch background pixel
			color := paletteColors[colorId]

			//TODO fetch plane pixel
			pixel, palettId := ppu.shiftRegiseter.PopPixel()
			if pixel != 0 {
				colorId := ppu.readRam(uint16(0x3F00) + uint16(palettId+pixel))
				color = paletteColors[colorId]
			}

			//TODO fetch sprite pixel

			if ppu.screenPosX < 256 && ppu.screenPosY < 240 {
				ppu.screen.setDot(ppu.screenPosX, ppu.screenPosY, color)
			} else {
				fmt.Printf("Wrong screen coordinates %d %d\n", ppu.screenPosX, ppu.screenPosY)
			}
		}
		if (line == Visible || line == PreRender) && (dotNum < 256 || (dotNum >= 321 && dotNum < 338)) {
			ppu.updatePlaneDataBuffer(dotNum)
			if dotNum%8 == 0 {
				dataLow := ppu.planeDataBuffer.tileDataLow
				dataHigh := ppu.planeDataBuffer.tileDataHigh
				paletteId := ppu.planeDataBuffer.attributeData
				ppu.shiftRegiseter.PushData(dataLow, dataHigh, paletteId)
			}
		}

		if line == Visible && dotNum < 256 {
			ppu.internalAddrReg.IncrementCoarseX()
			ppu.screenPosX += 1
		}
		if line == Visible && dotNum == 256 {
			ppu.internalAddrReg.IncrementY()
			ppu.screenPosY += 1
			ppu.screenPosX = 0
		}
		if dotNum == 257 { // copy horizontal position from tmp register
			ppu.internalAddrReg.CopyHorizontalPosition()
		}
		if line == PreRender && (dotNum >= 280 && dotNum < 304) { // copy vertical position from tmp register
			ppu.internalAddrReg.CopyVerticalPosition()
		}
	} else {
		colorId := rand.Intn(len(paletteColors))
		color := paletteColors[colorId]
		ppu.screen.setDot(ppu.screenPosX, ppu.screenPosY, color)
		ppu.screenPosX += 1
		if ppu.screenPosX == 256 {
			ppu.screenPosY += 1
			ppu.screenPosX = 0
		}
		if ppu.screenPosY == 240 {
			ppu.screenPosY = 0
		}
	}

	if line == Visible || line == PreRender && dotNum >= 257 && dotNum <= 320 {
		ppu.oamAddressReg.Reset()
	}
	if lineNum == 241 && dotNum == 1 {
		ppu.statusReg.SetStatusFlag(reg.V, true)
		if ppu.controllReg.GenerateNmi() {
			ppu.interruptSignal = true
		}
	}
	if line == PreRender && dotNum == 1 {
		ppu.statusReg.SetStatusFlag(reg.V, false)
		ppu.statusReg.SetStatusFlag(reg.S, false)
		ppu.statusReg.SetStatusFlag(reg.O, false)
		ppu.screenPosX = 0
		ppu.screenPosY = 0
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

func (ppu *Ppu) writeRam(address uint16, data byte) {
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

func (ppu *Ppu) updatePlaneDataBuffer(dotNum int) {
	switch dotNum % 8 {
	case 1:
		tileAddress := 0x2000 | ppu.internalAddrReg.GetAddress()&0x0FFF
		nameTableByte := ppu.readRam(tileAddress)
		ppu.planeDataBuffer.setTileId(nameTableByte)
	case 3:
		scrollAddress := ppu.internalAddrReg.GetAddress()
		attributeAddress := 0x23C0 | scrollAddress&0x0C00 | (scrollAddress>>4)&0x38 | (scrollAddress>>2)&0x7
		attributeByte := ppu.readRam(attributeAddress)
		if ppu.internalAddrReg.GetCoarseX()%2 == 1 {
			attributeByte >>= 2
		}
		if ppu.internalAddrReg.GetCoarseY()%2 == 1 {
			attributeByte >>= 2
		}
		attributeData := attributeByte & 0x3
		ppu.planeDataBuffer.setAttributeData(attributeData)
	case 5:
		tileAddress := uint16(ppu.planeDataBuffer.taileId) + ppu.internalAddrReg.GetFineY()
		data := ppu.readRam(tileAddress)
		ppu.planeDataBuffer.setTileDataLow(data)
	case 7:
		tileAddress := uint16(ppu.planeDataBuffer.taileId) + ppu.internalAddrReg.GetFineY() + 8
		data := ppu.readRam(tileAddress)
		ppu.planeDataBuffer.setTileDataHigh(data)
	}
}
