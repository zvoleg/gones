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

	var paletteWord uint32
	switch paletteId {
	case 0b00:
		paletteWord = 0x0000
	case 0b01:
		paletteWord = 0x5555 // 0b01010101....
	case 0b10:
		paletteWord = 0xAAAA // 0b10101010....
	case 0b11:
		paletteWord = 0xFFFF // 0b11111111....
	}
	paletteWord <<= fineX * 2
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
