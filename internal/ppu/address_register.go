package ppu

type addressRegister struct {
	value uint16
	latch *bool
}

func (r *addressRegister) write(value byte) {
	if !(*r.latch) {
		r.value = r.value & 0x00FF
		r.value = r.value | uint16(value)<<8
	} else {
		r.value = r.value & 0xFF00
		r.value = r.value | uint16(value)
	}
	*r.latch = !(*r.latch)
}

func (r *addressRegister) increment(incrementer uint16) {
	r.value += incrementer
}
