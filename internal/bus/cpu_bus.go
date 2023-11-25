package bus

func (bus *Bus) CpuRead(address uint16) byte {
	var data byte = 0
	if address <= 0x1FFF {
		data = bus.ram[address&0x07FF]
	} else if address >= 0x2000 && address <= 0x3FFF { // PPU registers
		// address & 0x7
	} else if address >= 0x4020 {
		data = bus.rom.ReadPrgRom(address)
	}
	return data
}

func (bus *Bus) CpuWrite(address uint16, data byte) {
	if address <= 0x1FFF {
		bus.ram[address&0x07FFF] = data
	} else if address >= 0x4020 {
		bus.rom.WritePrgRom(address, data)
	}
}
