package ppu

import (
	reg "github.com/zvoleg/gones/internal/ppu/internal/registers"
)

func (ppu *Ppu) RegisterRead(regAddress uint16) byte {
	var data byte
	switch regAddress {
	case 2:
		data = ppu.statusReg.Read()
		ppu.internalAddrReg.ResetLatch()
		ppu.statusReg.SetStatusFlag(reg.V, false)
	case 7:
		data = ppu.dataBuffer
		ppu.dataBuffer = ppu.readRam(ppu.internalAddrReg.GetAddress())
		if ppu.internalAddrReg.GetAddress() >= 0x3F00 {
			data = ppu.dataBuffer
		}
		ppu.internalAddrReg.Increment(ppu.controllReg.Incrementer())
	}
	return data
}

func (ppu *Ppu) RegisterWrite(regAddress uint16, data byte) {
	switch regAddress {
	case 0:
		ppu.controllReg.Write(data)
		ppu.internalAddrReg.SetNameTable(data)
	case 1:
		ppu.maskReg.Write(data)
	case 3:
		ppu.oamAddressReg.Write(data)
	case 4:
		// fmt.Printf("Write into SRAM, address: %04X\n", ppu.oamAddressReg.value)
		ppu.sram[ppu.oamAddressReg.GetAddress()] = data
		ppu.oamAddressReg.Increment()
	case 5:
		ppu.internalAddrReg.ScrollWrite(data)
		// fmt.Printf("Scroll set intern: 0x%04X\n", ppu.internalAddrReg.GetAddress())
	case 6:
		ppu.internalAddrReg.AddressWrite(data)
		// fmt.Printf("Address set intern: 0x%04X\n", ppu.internalAddrReg.GetAddress())
	case 7:
		ppu.writeRam(ppu.internalAddrReg.GetAddress(), data)
		ppu.internalAddrReg.Increment(ppu.controllReg.Incrementer())
	}
}
