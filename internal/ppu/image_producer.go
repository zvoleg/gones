package ppu

func (i *image) setDot(x, y int, dot color) {
	address := (y * i.width * 4) + x*4
	i.buff[address] = dot.r
	i.buff[address+1] = dot.g
	i.buff[address+2] = dot.b
	i.buff[address+3] = dot.a
}

func (ppu *Ppu) GetMainScreen() []byte {
	return ppu.screen.buff
}

func (ppu *Ppu) GetPatternTables() []byte {
	img := newImage(256, 128)
	ppu.readPatternTable(&img, 0)
	ppu.readPatternTable(&img, 1)
	return img.buff
}

func (ppu *Ppu) readPatternTable(img *image, table int) {
	posX := 128 * table
	posY := 0
	startAddress := 0x1000 * uint16(table)
	for tileAddress := startAddress; tileAddress < startAddress+0x1000; tileAddress += 0x10 {
		ppu.drawTile(img, tileAddress, posX, posY)
		posX += 8
		if posX >= 128+128*table {
			posX = 128 * table
			posY += 8
		}
	}
}

func (ppu *Ppu) GetNameTable() []byte {
	img := newImage(512, 512)
	ppu.readNameTable(&img, 0)
	ppu.readNameTable(&img, 1)
	ppu.readNameTable(&img, 2)
	ppu.readNameTable(&img, 3)
	return img.buff
}

func (ppu *Ppu) readNameTable(img *image, table int) {
	posX := 256 * (table % 2)
	posY := 256 * (table / 2)
	startAddress := 0x400*uint16(table) + 0x2000
	pattenrsTableAddr := ppu.controllReg.GetBackgroundTable()
	for i := startAddress; i < startAddress+0x400; i += 1 {
		tileId := uint16(ppu.readRam(i))
		tileAddress := pattenrsTableAddr + tileId*0x10
		ppu.drawTile(img, tileAddress, posX, posY)
		posX += 8
		if posX >= 256+256*(table%2) {
			posX = 256 * (table % 2)
			posY += 8
		}
	}
}

func (ppu *Ppu) GetColorPalette() []byte {
	img := newImage(9, 5)
	x := 0
	y := 0
	color := ppu.getColor(0, 0)
	img.setDot(x, y, color)
	y += 1
	for paletteId := 0; paletteId < 4; paletteId += 1 {
		for colorId := 0; colorId < 4; colorId += 1 {
			if paletteId == 0 && colorId == 0 {
				continue
			}
			color = ppu.getColor(byte(paletteId), byte(colorId))
			img.setDot(x, y, color)
			x += 1
			if x == 4 {
				x = 0
				y += 1
			}
		}
	}
	y = 1
	x = 5
	for paletteId := 4; paletteId < 8; paletteId += 1 {
		for colorId := 0; colorId < 4; colorId += 1 {
			if paletteId == 4 && colorId == 0 {
				continue
			}
			color = ppu.getColor(byte(paletteId), byte(colorId))
			img.setDot(x, y, color)
			x += 1
			if x == 9 {
				x = 5
				y += 1
			}
		}
	}
	return img.buff
}

func (ppu *Ppu) getColor(paletteId byte, dotBits byte) color {
	var colorId byte
	if dotBits != 0 {
		colorId = ppu.paletteRam[paletteId*4+dotBits]
	} else {
		colorId = ppu.paletteRam[0]
	}
	return paletteColors[colorId]
}

func (ppu *Ppu) drawTile(img *image, tileAddress uint16, posX, posY int) {
	for tileLineNum := 0; tileLineNum < 8; tileLineNum += 1 {
		lineAddress := tileAddress + uint16(tileLineNum)
		planeLow := ppu.readRam(lineAddress)
		planeHigh := ppu.readRam(lineAddress + 8)
		for dotNum := 0; dotNum < 8; dotNum += 1 {
			offset := 7 - dotNum
			dotBits := ((planeHigh>>offset)<<1)&2 | (planeLow>>offset)&1
			d := ppu.getColor(4, dotBits)
			img.setDot(posX+dotNum, posY+tileLineNum, d)
		}
	}
}
