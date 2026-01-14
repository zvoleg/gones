package ppu

func (ppu *Ppu) InitDma(page byte) {
	ppu.dmaSrcPage = uint16(page) << 8
	ppu.dmaEnabled = true
	ppu.dmaClockWaiter = true
}

func (ppu *Ppu) DmaClockRead() {
	ppu.dmaByte = ppu.bus.ReadDmaByte(ppu.dmaSrcPage)
	ppu.dmaSrcPage += 1
}

func (ppu *Ppu) DmaClockWrite() {
	sramAddress := ppu.oamAddressReg.GetAddress()
	ppu.writeSram(ppu.dmaByte, sramAddress)
	ppu.oamAddressReg.Increment()
	if ppu.dmaSrcPage&0x00FF == 0 {
		ppu.dmaEnabled = false
	}
}

func (ppu *Ppu) DmaEnable() bool {
	return ppu.dmaEnabled
}
