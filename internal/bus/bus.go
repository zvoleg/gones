package bus

type Cartridge interface {
	ReadPrgRom(address uint16) byte
	ReadChrRom(address uint16) byte
	WritePrgRom(address uint16, data byte)
	WriteChrRom(address uint16, data byte)
}

type PpuRegisterMap interface {
	RegisterRead(regAddress uint16) byte
	RegisterWrite(regAddress uint16, data byte)
}

type Bus struct {
	ram            [0x0800]byte
	rom            Cartridge
	ppuRegisterMap PpuRegisterMap
}

func New(rom Cartridge, ppuRegisterMap PpuRegisterMap) Bus {
	return Bus{
		rom:            rom,
		ppuRegisterMap: ppuRegisterMap,
	}
}
