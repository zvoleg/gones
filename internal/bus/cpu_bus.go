package bus

func (bus *Bus) CpuRead(address uint16) uint8 {
	return bus.ram[address]
}

func (bus *Bus) CpuWrite(address uint16, data uint8) {
	bus.ram[address] = data
}
