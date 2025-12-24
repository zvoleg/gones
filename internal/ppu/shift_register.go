package ppu

type shiftRegiseter struct {
	valueLow  uint16
	valueHigh uint16
	paletteId uint32
}

func (reg *shiftRegiseter) pushData(low, high, paletteId byte) {
	reg.valueLow &= 0xFF00 // clear low byte before setting
	reg.valueLow |= uint16(low)
	reg.valueHigh &= 0xFF00
	reg.valueHigh |= uint16(high)

	paletteByte := paletteId<<6 | paletteId<<4 | paletteId<<2 | paletteId
	reg.paletteId |= uint32(paletteByte)
}

func (reg *shiftRegiseter) scrollX(x int) {
	reg.valueLow <<= x
	reg.valueHigh <<= x
	reg.paletteId <<= x * 2
}

func (reg *shiftRegiseter) popPixel() (byte, byte) {
	lowBit := (reg.valueLow) >> 15 & 1
	highBit := (reg.valueHigh >> 15) & 1
	pixel := (highBit << 1) | lowBit
	palletId := (reg.paletteId >> 30) & 0x3
	reg.scrollX(1)
	return byte(pixel), byte(palletId)
}
