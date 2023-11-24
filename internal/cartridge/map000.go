package cartridge

type Map000 struct {
	header header
}

func newMap000(header header) Map000 {
	return Map000{header}
}

func (m *Map000) prgAddress(address uint16) uint16 {
	var mappedAddress uint16
	switch m.header.prgRomUnits {
	case 1:
		mappedAddress = address & 0x3FFF
	case 2:
		mappedAddress = address & 0x7FFF
	}
	return mappedAddress
}

func (m *Map000) chrAddress(address uint16) uint16 {
	return 0
}
