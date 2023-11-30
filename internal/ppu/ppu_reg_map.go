package ppu

import "fmt"

func (ppu *Ppu) RegisterRead(regAddress uint16) byte {
	var data byte
	switch regAddress {
	case 2:
		data = ppu.statusReg.read()
	}
	return data
}

func (ppu *Ppu) RegisterWrite(regAddress uint16, data byte) {
	switch regAddress {
	case 0:
		ppu.controllReg.write(data)
	case 1:
		ppu.maskReg.write(data)
	case 3:
		ppu.oamAddressReg.write(data)
	case 4:
		fmt.Println("Write into SRAM, address: ", ppu.oamAddressReg.value)
		ppu.sram[ppu.oamAddressReg.value] = data
		ppu.oamAddressReg.increment()
	case 5:
		ppu.scrollRegister.write(data)
	case 6:
		ppu.addressRegister.write(data)
	case 7:
		address := ppu.addressRegister.value
		fmt.Println("Write into VRAM, address: ", address)
		switch true {
		case address <= 0x1FFF:
			ppu.patternTable[address] = data
		case address >= 0x2000 && address <= 0x3EFF:
			address = address & 0x1FFF
			ppu.nameTable[address] = data
		case address >= 0x3F00 && address <= 0x3FFF:
			address = address & 0x1F
			ppu.paletteRam[address] = data
		default:
			fmt.Println("Wrong address for writing into vram: ", address)
		}
		incrementer := uint16(ppu.controllReg.getVramIncrement())
		ppu.addressRegister.increment(incrementer)
	}
}
