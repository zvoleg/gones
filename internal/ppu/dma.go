package ppu

func (ppu *Ppu) InitDma(page byte) {
	ppu.dmaEnabled = true
	ppu.dmaClockWaiter = true
}

func (ppu *Ppu) DmaClock() {
	if ppu.clockCounter%2 == 0 {
		ppu.dmaByte = ppu.bus.ReadDmaByte(ppu.oamAddressReg.GetAddress())
	} else {
		sramAddress := ppu.oamAddressReg.GetAddress() & 0x00FF
		ppu.sram[sramAddress] = ppu.dmaByte
		ppu.oamAddressReg.Increment()
		if ppu.oamAddressReg.GetAddress()&0x00FF == 0 {
			ppu.dmaEnabled = false
		}
	}
}

func (ppu *Ppu) DmaEnable() bool {
	return ppu.dmaEnabled
}
