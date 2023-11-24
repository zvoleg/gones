package bus

type Cartridge interface {
	ReadPrgRom(address uint16) uint8
	ReadChrRom(address uint16) uint8
	WritePrgRom(address uint16, data uint8)
	WriteChrRom(address uint16, data uint8)
}

type Bus struct {
	ram [0x0800]uint8
	rom Cartridge
}

func New(rom Cartridge) Bus {
	return Bus{rom: rom}
}
