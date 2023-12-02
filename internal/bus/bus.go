package bus

type Cartridge interface {
	ReadPrgRom(address uint16) byte
	ReadChrRom(address uint16) byte
	WritePrgRom(address uint16, data byte)
	WriteChrRom(address uint16, data byte)
}

type PpuExposedInterface interface {
	RegisterRead(regAddress uint16) byte
	RegisterWrite(regAddress uint16, data byte)
	InitDma(page byte)
}

type Bus struct {
	ram            [0x0800]byte
	rom            Cartridge
	ppuRegisterMap PpuExposedInterface
}

func New(rom Cartridge, ppuRegisterMap PpuExposedInterface) Bus {
	return Bus{
		rom:            rom,
		ppuRegisterMap: ppuRegisterMap,
	}
}
