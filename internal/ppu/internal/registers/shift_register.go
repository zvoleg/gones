package registers

type ShiftRegiseter struct {
	valueLow  uint16
	valueHigh uint16
	paletteId uint32
}

func (reg *ShiftRegiseter) PushData(low, high, paletteId byte) {
	reg.valueLow &= 0xFF00 // clear low byte before setting
	reg.valueLow |= uint16(low)
	reg.valueHigh &= 0xFF00
	reg.valueHigh |= uint16(high)

	paletteByte := paletteId<<6 | paletteId<<4 | paletteId<<2 | paletteId
	reg.paletteId |= uint32(paletteByte)
}

func (reg *ShiftRegiseter) ScrollX(x int) {
	reg.valueLow <<= x
	reg.valueHigh <<= x
	reg.paletteId <<= x * 2
}

func (reg *ShiftRegiseter) PopPixel() (byte, byte) {
	lowBit := (reg.valueLow) >> 15 & 1
	highBit := (reg.valueHigh >> 15) & 1
	pixel := (highBit << 1) | lowBit
	palletId := (reg.paletteId >> 30) & 0x3
	reg.ScrollX(1)
	return byte(pixel), byte(palletId)
}
