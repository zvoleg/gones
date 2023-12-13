package ppu

import "fmt"

func (ppu *Ppu) RegisterRead(regAddress uint16) byte {
	var data byte
	switch regAddress {
	case 2:
		data = ppu.statusReg.read()
		*ppu.addressRegister.latch = false
		ppu.statusReg.setStatusFlag(V, false)
	case 7:
		data = ppu.dataBuffer
		ppu.dataBuffer = ppu.readVram()
		if ppu.addressRegister.value >= 0x3F00 {
			data = ppu.dataBuffer
		}
		ppu.addressRegister.increment(ppu.controllReg.incrementer)
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
		fmt.Printf("Write into SRAM, address: %04X\n", ppu.oamAddressReg.value)
		ppu.sram[ppu.oamAddressReg.value] = data
		ppu.oamAddressReg.increment()
	case 5:
		ppu.scrollRegister.write(data)
	case 6:
		ppu.addressRegister.write(data)
	case 7:
		ppu.writeVram(data)
		ppu.addressRegister.increment(ppu.controllReg.incrementer)
	}
}
