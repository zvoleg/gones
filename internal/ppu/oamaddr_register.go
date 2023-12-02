package ppu

type oamAddressReg struct {
	value uint16
}

func (r *oamAddressReg) write(value byte) {
	r.value = uint16(value) << 8
}

func (r *oamAddressReg) increment() {
	r.value += 1
}
