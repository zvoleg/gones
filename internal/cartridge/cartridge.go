package cartridge

import (
	"os"

	"github.com/zvoleg/gones/internal/ppu"
)

const prgRomUnitSize = 0x4000
const chrRomUnitSize = 0x2000

type mapper interface {
	prgAddress(address uint16) uint16
	chrAddress(address uint16) uint16
}

type Cartridge struct {
	header header
	mapper mapper
	prgRom []byte
	chrRom []byte
}

func New(filePath string) Cartridge {
	nesFile, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	header := new(nesFile[0:16])
	mapper := getMapper(header, nesFile)
	prgRomStartOffset := 16
	if header.trainer {
		prgRomStartOffset += 512
	}
	prgRomEndOffset := prgRomStartOffset + header.prgRomUnits*prgRomUnitSize
	prgRom := nesFile[prgRomStartOffset:prgRomEndOffset]
	chrRomStartOffset := prgRomEndOffset
	chrRomEndOffset := chrRomStartOffset + header.chrRomUnits*chrRomUnitSize
	chrRom := nesFile[chrRomStartOffset:chrRomEndOffset]

	return Cartridge{
		header: header,
		mapper: mapper,
		prgRom: prgRom,
		chrRom: chrRom,
	}
}

func getMapper(header header, nesFile []byte) mapper {
	var mapp mapper
	switch header.mapperNum {
	case 000:
		map000 := newMap000(header)
		mapp = &map000
	}
	return mapp
}

func (c *Cartridge) ReadPrgRom(address uint16) byte {
	mappedAddress := c.mapper.prgAddress(address)
	return c.prgRom[mappedAddress]
}

func (c *Cartridge) ReadChrRom(address uint16) byte {
	mappedAddress := c.mapper.chrAddress(address)
	return c.chrRom[mappedAddress]
}

func (c *Cartridge) WritePrgRom(address uint16, data byte) {
	// mappedAddress := c.mapper.prgAddress(address)
}

func (c *Cartridge) WriteChrRom(address uint16, data byte) {
	// mappedAddress := c.mapper.prgAddress(address)
}

func (c *Cartridge) Mirroring() ppu.Mirroring {
	return c.header.mirroring
}
