package bus

import "github.com/zvoleg/gones/internal/ppu"

func (bus *Bus) PpuRead(address uint16) byte {
	return bus.rom.ReadChrRom(address)
}

func (bus *Bus) PpuWrite(address uint16, data byte) {

}

func (bus *Bus) ReadDmaByte(address uint16) byte {
	return bus.CpuRead(address)
}

func (bus *Bus) GetMirroring() ppu.Mirroring {
	return bus.rom.Mirroring()
}
