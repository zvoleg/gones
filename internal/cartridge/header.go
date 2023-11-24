package cartridge

import "github.com/zvoleg/gones/internal/ppu"

type header struct {
	prgRomUnits int
	chrRomUnits int
	mapperNum   int
	mirroring   ppu.Mirroring
	battery     bool
	trainer     bool
	fourScreen  bool
}

// implementation only for the iNes format (iNes2.0 not supported)
func new(rawHeader []byte) header {
	fileFormat := string(rawHeader[0:3])
	if fileFormat != "NES" {
		panic("Unexpected format bytes" + fileFormat)
	}
	prgRomUnits := int(rawHeader[4])
	chrRomUnits := int(rawHeader[5])

	flag6 := rawHeader[6]
	mapperNumL := flag6 >> 4
	mirroring := ppu.Mirroring(flag6 & 0x1)
	battery := flag6&0x02 == 1
	trainer := flag6&0x04 == 1
	fourScreen := flag6&0x08 == 1

	flag7 := rawHeader[7]
	mapperNumH := flag7 >> 4
	iNes2 := flag7&0x0C == 0x0C
	if iNes2 {
		panic("iNes2.0 format is not supported")
	}

	mapperNum := int(mapperNumH<<4 | mapperNumL)

	return header{
		prgRomUnits: prgRomUnits,
		chrRomUnits: chrRomUnits,
		mapperNum:   mapperNum,
		mirroring:   mirroring,
		battery:     battery,
		trainer:     trainer,
		fourScreen:  fourScreen,
	}
}
