package bus

type Cartridge interface {
	ReadPrgRom(address uint16) byte
	ReadChrRom(address uint16) byte
	WritePrgRom(address uint16, data byte)
	WriteChrRom(address uint16, data byte)
}

type Bus struct {
	ram [0x0800]byte
	rom Cartridge
}

func New(rom Cartridge) Bus {
	return Bus{rom: rom}
}
