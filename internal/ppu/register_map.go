package ppu

import "fmt"

func (ppu *Ppu) RegisterRead(regAddress uint16) byte {
	var data byte
	switch regAddress {
	case 2:
		data = ppu.statusReg.read()
		ppu.internalAddrReg.resetLatch()
		ppu.statusReg.setStatusFlag(V, false)
	case 7:
		data = ppu.dataBuffer
		ppu.dataBuffer = ppu.readRam(ppu.internalAddrReg.cur_value)
		if ppu.internalAddrReg.cur_value >= 0x3F00 {
			data = ppu.dataBuffer
		}
		ppu.internalAddrReg.increment(ppu.controllReg.incrementer)
		ppu.addressRegister.increment(ppu.controllReg.incrementer)
	}
	return data
}

func (ppu *Ppu) RegisterWrite(regAddress uint16, data byte) {
	switch regAddress {
	case 0:
		ppu.controllReg.write(data)
		ppu.internalAddrReg.setNameTable(data)
	case 1:
		ppu.maskReg.write(data)
	case 3:
		ppu.oamAddressReg.write(data)
	case 4:
		// fmt.Printf("Write into SRAM, address: %04X\n", ppu.oamAddressReg.value)
		ppu.sram[ppu.oamAddressReg.value] = data
		ppu.oamAddressReg.increment()
	case 5:
		ppu.internalAddrReg.scrollWrite(data)
	case 6:
		ppu.internalAddrReg.addressWrite(data)
		ppu.addressRegister.write(data)
		fmt.Printf("addr: 0x%04X | intern: 0x%04X\n", ppu.addressRegister.value, ppu.internalAddrReg.cur_value)
	case 7:
		ppu.writeVram(ppu.internalAddrReg.cur_value, data)
		ppu.internalAddrReg.increment(ppu.controllReg.incrementer)
		ppu.addressRegister.increment(ppu.controllReg.incrementer)
	}
}
