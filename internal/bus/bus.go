package bus

import (
	"github.com/zvoleg/gones/internal/controller"
	"github.com/zvoleg/gones/internal/ppu"
)

type Cartridge interface {
	ReadPrgRom(address uint16) byte
	ReadChrRom(address uint16) byte
	WritePrgRom(address uint16, data byte)
	WriteChrRom(address uint16, data byte)
	Mirroring() ppu.Mirroring
}

type PpuRegisters interface {
	RegisterRead(regAddress uint16) byte
	RegisterWrite(regAddress uint16, data byte)
	InitDma(page byte)
}

type Bus struct {
	ram            [0x0800]byte
	rom            Cartridge
	ppuRegisterMap PpuRegisters
	inputInterface controller.InputInterface
}

func New(rom Cartridge, ppuRegisterMap PpuRegisters, inputInterface controller.InputInterface) Bus {
	return Bus{
		rom:            rom,
		ppuRegisterMap: ppuRegisterMap,
		inputInterface: inputInterface,
	}
}
