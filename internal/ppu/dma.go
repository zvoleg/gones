package ppu

import "github.com/zvoleg/gones/internal/cpu6502"

func (ppu *Ppu) InitDma(page byte) {
	ppu.dmaEnabled = true
	ppu.cpuSignalLine <- cpu6502.DmaEnable
}

func (ppu *Ppu) DmaClock() {
	if ppu.clockCounter%2 == 0 {
		ppu.bus.ReadDmaByte(ppu.oamAddressReg.value)
	} else {
		sramAddress := ppu.oamAddressReg.value & 0x00FF
		ppu.sram[sramAddress] = 0
		ppu.oamAddressReg.increment()
		if ppu.oamAddressReg.value&0x00FF == 0 {
			ppu.dmaEnabled = false
			ppu.cpuSignalLine <- cpu6502.DmaDisable
		}
	}
}
