package ppu

import (
	"fmt"
	"math/bits"
	"math/rand"
	"unsafe"

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

	oam                [0x40]objectAttributeEntity // sprites memory
	nextLineOam        [0x08]objectAttributeEntity
	curentLineOam      [0x08]objectAttributeEntity
	oamCounter         int
	nextLineOamCounter int

	screen image // screen size 256x240, 4 byte per pixel (RGBA)

	planeDataBuffer planeDataBuffer
	shiftRegiseter  reg.ShiftRegiseter

	controllReg     *reg.ControllReg
	maskReg         *reg.MaskReg
	statusReg       *reg.StatusReg
	oamAddressReg   *reg.OamAddressReg
	internalAddrReg *reg.InternalAddrReg

	dataBuffer byte

	bus PpuBus

	dmaEnabled     bool
	dmaSrcPage     uint16
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
	internalAddrReg := reg.InternalAddrReg{}

	return Ppu{
		screen: newImage(256, 240),

		controllReg:     &controllReg,
		maskReg:         &maskReg,
		statusReg:       &statusReg,
		oamAddressReg:   &oamAddressReg,
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
	// if ppu.clockCounter == 0 && !ppu.evenFrame {
	// 	ppu.clockCounter += 1
	// 	return
	// }

	lineNum := ppu.clockCounter / 341
	dotNum := ppu.clockCounter % 341
	line := getLineType(lineNum)

	if ppu.maskReg.RenderingEnabled() {
		if line == Visible || line == PreRender {
			ppu.backgroundPlaneProcess(dotNum, lineNum)
		}
		if line == PreRender && (dotNum >= 280 && dotNum < 305) { // copy vertical position from tmp register
			ppu.internalAddrReg.CopyVerticalPosition()
		}
		if line == Visible && dotNum <= 257 {
			ppu.spriteEvaluations(dotNum, lineNum)
		}
		if line == Visible && dotNum < 256 {
			ppu.rendering(dotNum, lineNum)
		}
	} else {
		colorId := rand.Intn(len(paletteColors))
		color := paletteColors[colorId]
		if dotNum < 256 && lineNum < 240 {
			ppu.screen.setDot(dotNum, lineNum, color)
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
	address &= 0x3FFF
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
		fmt.Printf("Wrong address for reading into vram: 0x%04X\n", address)
	}
	return data
}

func (ppu *Ppu) writeRam(address uint16, data byte) {
	// fmt.Printf("Write into VRAM, address: %04X\n", address)
	address &= 0x3FFF
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

func (ppu *Ppu) writeSram(data byte, address byte) {
	ptr := unsafe.Pointer(&ppu.oam[0])
	offset := unsafe.Pointer(uintptr(ptr) + uintptr(address))
	bytePtr := (*byte)(offset)
	*bytePtr = data
}

func (ppu *Ppu) readSram(address byte) byte {
	ptr := unsafe.Pointer(&ppu.oam[0])
	offset := unsafe.Pointer(uintptr(ptr) + uintptr(address))
	bytePtr := (*byte)(offset)
	return *bytePtr
}

func (ppu *Ppu) readTileLine(table, tileId, tileLine uint16) (byte, byte) {
	tileDataLow := ppu.readRam(table + tileId*16 + tileLine)
	tileDataHigh := ppu.readRam(table + tileId*16 + 8 + tileLine)
	return tileDataLow, tileDataHigh
}

func (ppu *Ppu) updatePlaneDataBuffer(dotNum int) {
	switch dotNum % 8 {
	case 0:
		dataLow := ppu.planeDataBuffer.tileDataLow
		dataHigh := ppu.planeDataBuffer.tileDataHigh
		paletteId := ppu.planeDataBuffer.attributeData
		ppu.shiftRegiseter.PushData(dataLow, dataHigh, paletteId, ppu.internalAddrReg.GetFineX())
		ppu.internalAddrReg.IncrementCoarseX()
	case 1:
		tileAddress := 0x2000 | ppu.internalAddrReg.GetAddress()&0x0FFF
		// fmt.Printf("tile address: %04X\n", tileAddress)
		tileId := ppu.readRam(tileAddress)
		ppu.planeDataBuffer.setTileId(tileId)
	case 3:
		scrollAddress := ppu.internalAddrReg.GetAddress()
		coarseY := ppu.internalAddrReg.GetCoarseY()
		coarseX := ppu.internalAddrReg.GetCoarseX()
		attributeAddress := 0x23C0 | scrollAddress&0x0C00 | (coarseY>>2)<<3 | (coarseX >> 2)
		attributeByte := ppu.readRam(attributeAddress)
		if coarseX&0x02 == 0x02 {
			attributeByte >>= 2
		}
		if coarseY&0x02 == 0x02 {
			attributeByte >>= 4
		}
		attributeData := attributeByte & 0x03
		ppu.planeDataBuffer.setAttributeData(attributeData)
	case 5:
		backgroundTable := ppu.controllReg.GetBackgroundTable()
		tileAddress := backgroundTable + uint16(ppu.planeDataBuffer.taileId)*16 + ppu.internalAddrReg.GetFineY()
		data := ppu.readRam(tileAddress)
		ppu.planeDataBuffer.setTileDataLow(data)
	case 7:
		backgroundTable := ppu.controllReg.GetBackgroundTable()
		tileAddress := backgroundTable + uint16(ppu.planeDataBuffer.taileId)*16 + ppu.internalAddrReg.GetFineY() + 8
		data := ppu.readRam(tileAddress)
		ppu.planeDataBuffer.setTileDataHigh(data)
	}
}

func (ppu *Ppu) checkCurrentLineOam(dotNum int) []objectAttributeEntity {
	hitObjects := make([]objectAttributeEntity, 0, 8)
	for i, entity := range ppu.curentLineOam {
		xDiff := dotNum - int(entity.xCoordinate)
		if xDiff >= 0 && xDiff < 8 {
			hitObjects = append(hitObjects, ppu.curentLineOam[i])
		}
	}
	return hitObjects
}

func (ppu *Ppu) backgroundPlaneProcess(dotNum, lineNum int) {
	if (dotNum > 0 && dotNum < 256) || (dotNum >= 321 && dotNum < 337) {
		ppu.shiftRegiseter.ScrollX(1)
		ppu.updatePlaneDataBuffer(dotNum)
	}
	if dotNum == 256 {
		ppu.internalAddrReg.IncrementY()
	}
	if dotNum == 257 { // copy horizontal position from tmp register
		ppu.internalAddrReg.CopyHorizontalPosition()
	}
}

func (ppu *Ppu) spriteEvaluations(dotNum, lineNum int) {
	if dotNum == 1 {
		for i := range 8 {
			ppu.nextLineOam[i].toDefault()
		}
	}
	if dotNum >= 65 && dotNum < 257 && ppu.oamCounter < 64 && ppu.nextLineOamCounter < 9 {
		candidate := ppu.oam[ppu.oamCounter]
		ppu.oamCounter += 1
		yOffset := lineNum + 1 - int(candidate.yCoordinate)
		if yOffset >= 0 && yOffset < int(ppu.controllReg.SpriteSize()) {
			if ppu.nextLineOamCounter < 8 {
				ppu.nextLineOam[ppu.nextLineOamCounter] = candidate
			} else {
				ppu.statusReg.SetStatusFlag(reg.O, true)
			}
			ppu.nextLineOamCounter += 1
		}
	}
	if dotNum == 257 {
		ppu.curentLineOam, ppu.nextLineOam = ppu.nextLineOam, ppu.curentLineOam
		ppu.oamCounter = 0
		ppu.nextLineOamCounter = 0
	}
}

func (ppu *Ppu) rendering(dotNum, lineNum int) {
	// fmt.Printf("Dot: %d\tLine: %d\tInternal Address: 0x%04X\n", dotNum, lineNum, ppu.internalAddrReg.GetAddress())
	colorId := ppu.readRam(0x3F00) // fetch background pixel
	color := paletteColors[colorId]

	// fetch plane pixel
	backgroundPixel := 0
	if ppu.maskReg.LeftBackgroundEnabled() || dotNum >= 8 {
		backgroundPixel, palettId := ppu.shiftRegiseter.PopPixel()
		if backgroundPixel != 0 {
			colorId := ppu.readRam(uint16(0x3F00) + uint16(palettId*4+backgroundPixel))
			color = paletteColors[colorId]
		}
	}

	// fetch sprite pixel
	if ppu.maskReg.LeftSpritesEnabled() || dotNum >= 8 {
		hitObjects := ppu.checkCurrentLineOam(dotNum)
		spriteTable := ppu.controllReg.GetSpriteTable()
		for _, entity := range hitObjects {
			tileId := entity.tileId
			tileLine := lineNum - int(entity.yCoordinate)
			if entity.VerticalFlip() {
				tileLine = 7 - tileLine
			}
			tileLow, tileHigh := ppu.readTileLine(spriteTable, uint16(tileId), uint16(tileLine))
			if entity.HorizontalFlip() {
				tileLow = bits.Reverse8(tileLow)
				tileHigh = bits.Reverse8(tileHigh)
			}

			xOffset := dotNum - int(entity.xCoordinate)
			spritePixel := ((tileHigh>>(7-byte(xOffset)))&0x1)<<1 | ((tileLow >> (7 - byte(xOffset))) & 0x1)
			if spritePixel == 0 {
				continue
			}
			if backgroundPixel != 0 && spritePixel != 0 {
				ppu.statusReg.SetStatusFlag(reg.S, true)
			}
			if entity.BehindBackground() && backgroundPixel != 0 {
				continue
			}
			colorId := ppu.readRam(uint16(0x3F00) + uint16(entity.Palette()*4+uint16(spritePixel)))
			color = paletteColors[colorId]
		}
	}

	ppu.screen.setDot(dotNum, lineNum, color)
}
