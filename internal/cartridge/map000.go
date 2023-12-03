package cartridge

import "fmt"

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
	if address < 0x2000 {
		return address
	}
	fmt.Println("Unexpected chrom address ", address)
	return 0
}
