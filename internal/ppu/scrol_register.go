package ppu

type scrollRegister struct {
	value uint16
	latch *bool
}

func (r *scrollRegister) write(value byte) {
	if !(*r.latch) {
		r.value = r.value & 0x00FF
		r.value = r.value | uint16(value)<<8
	} else {
		r.value = r.value & 0xFF00
		r.value = r.value | uint16(value)
	}
	*r.latch = !(*r.latch)
}
