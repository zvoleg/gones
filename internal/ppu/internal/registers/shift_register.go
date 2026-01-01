package registers

type ShiftRegiseter struct {
	valueLow  uint16
	valueHigh uint16
	paletteId uint32
}

func (reg *ShiftRegiseter) PushData(low, high, paletteId byte, fineX uint16) {
	reg.valueLow &= 0xFF00 << fineX // clear low byte before setting
	reg.valueLow |= uint16(low) << fineX
	reg.valueHigh &= 0xFF00 << fineX
	reg.valueHigh |= uint16(high) << fineX

	paletteBits := uint32(paletteId)
	paletteWord := (paletteBits<<14 | paletteBits<<12 | paletteBits<<10 | paletteBits<<8 | paletteBits<<6 | paletteBits<<4 | paletteBits<<2 | paletteBits) << (fineX * 2)
	reg.paletteId |= uint32(paletteWord)
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
	return byte(pixel), byte(palletId)
}
