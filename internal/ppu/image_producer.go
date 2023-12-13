package ppu

func (i *image) setDot(x, y int, dot color) {
	address := (y * i.width * 4) + x*4
	i.buff[address] = dot.r
	i.buff[address+1] = dot.g
	i.buff[address+2] = dot.b
	i.buff[address+3] = dot.a
}

func (ppu *Ppu) GetMainScreen() []byte {
	buffer := make([]byte, 1024)
	return buffer
}

func (ppu *Ppu) GetPatternTables() []byte {
	img := newImage(256, 128)
	ppu.readPatternTable(&img, 0)
	ppu.readPatternTable(&img, 1)
	return img.buff
}

func (ppu *Ppu) readPatternTable(img *image, table int) {
	startX := 128 * table
	x := startX
	xOffset := startX
	y := 0
	yOffset := 0
	startAddress := 0x1000 * table
	for i := startAddress; i < startAddress+0x1000; i += 0x10 {
		// tile line end
		if i != startAddress && i%0x100 == 0 {
			yOffset += 8
			xOffset = startX
			x = startX
			y = yOffset
		}
		for b := 0; b < 8; b += 1 {
			plane0 := ppu.bus.PpuRead(uint16(i + b))
			plane1 := ppu.bus.PpuRead(uint16(i + b + 8))
			for bit := 0; bit < 8; bit += 1 {
				offset := 7 - bit
				dotBits := ((plane1>>offset)<<1)&2 | (plane0>>byte(offset))&1
				var d color
				switch dotBits {
				case 0:
					d = newColor(0x00, 00, 00)
				case 1:
					d = newColor(0xD0, 0xD0, 0xD0)
				case 2:
					d = newColor(0x50, 0x50, 0xD0)
				case 3:
					d = newColor(0xD0, 0x50, 0x50)
				}
				img.setDot(x, y, d)
				x += 1
			}
			y += 1
			x = xOffset
		}
		// tile end
		xOffset += 8
		y = yOffset
		x = xOffset
	}
}

func (ppu *Ppu) GetNameTable() []byte {
	img := newImage(512, 512)
	ppu.readNameTable(&img, 0)
	ppu.readNameTable(&img, 1)
	return img.buff
}

func (ppu *Ppu) readNameTable(img *image, table int) {
	startX := 256 * table
	startY := 0
	startAddress := 0x400 * table
	for i := startAddress; i < startAddress+0x400; i += 1 {
		spiteId := uint16(ppu.nameTable[i])
		x := startX
		y := startY
		for spriteByteNum := 0; spriteByteNum < 8; spriteByteNum += 1 {
			spriteTable := ppu.controllReg.backgroundTable
			spriteAddress := spriteTable + spiteId*0x10 + uint16(spriteByteNum)
			plane0 := ppu.bus.PpuRead(spriteAddress)
			plane1 := ppu.bus.PpuRead(spriteAddress + 8)
			for bit := 0; bit < 8; bit += 1 {
				offset := 7 - bit
				dotBits := ((plane1>>offset)<<1)&2 | (plane0>>byte(offset))&1
				var d color
				switch dotBits {
				case 0:
					d = newColor(0x00, 00, 00)
				case 1:
					d = newColor(0xD0, 0xD0-0x50*byte(table), 0xD0-0x50*byte(table))
				case 2:
					d = newColor(0x50, 0x50, 0xD0)
				case 3:
					d = newColor(0xD0, 0x50, 0x50)
				}
				img.setDot(x, y, d)
				x += 1
			}
			x = startX
			y += 1
		}
		startX += 8
		if startX >= 256+256*table {
			startX = 256 * table
			startY += 8
		}
	}
}

func (ppu *Ppu) GetCollorPallete() []byte {
	img := newImage(9, 5)
	x := 0
	y := 0
	colorId := ppu.paletteRam[0]
	color := palletteColors[colorId]
	img.setDot(x, y, color)
	y += 1
	for _, colorId := range ppu.paletteRam[1:0x10] {
		color = palletteColors[colorId]
		img.setDot(x, y, color)
		x += 1
		if x == 4 {
			x = 0
			y += 1
		}
	}
	y = 0
	x = 5
	colorId = ppu.paletteRam[0x10]
	color = palletteColors[colorId]
	img.setDot(x, y, color)
	y += 1
	for _, colorId := range ppu.paletteRam[0x11:0x20] {
		color = palletteColors[colorId]
		img.setDot(x, y, color)
		x += 1
		if x == 9 {
			x = 5
			y += 1
		}
	}
	return img.buff
}
