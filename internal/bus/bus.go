package bus

import "github.com/zvoleg/gones/internal/ppu"

type Cartridge interface {
	ReadPrgRom(address uint16) byte
	ReadChrRom(address uint16) byte
	WritePrgRom(address uint16, data byte)
	WriteChrRom(address uint16, data byte)
}

type Bus struct {
	ram [0x0800]byte
	rom Cartridge

	// ppu interface
	controllReg   *ppu.ControllReg
	maskReg       *ppu.MaskReg
	statusReg     *ppu.StatusReg
	oamAddressReg *ppu.OamAddressReg
	oamDataReg    *ppu.OamDataReg
}

func New(rom Cartridge, ppu *ppu.Ppu) Bus {
	return Bus{
		rom:           rom,
		controllReg:   ppu.ControllReg,
		maskReg:       ppu.MaskReg,
		statusReg:     ppu.StatusReg,
		oamAddressReg: ppu.OamAddressReg,
		oamDataReg:    ppu.OamDataReg,
	}
}
