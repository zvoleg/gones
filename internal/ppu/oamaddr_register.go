package ppu

type oamAddressReg struct {
	value byte
}

func (r *oamAddressReg) write(value byte) {
	r.value = value
}

func (r *oamAddressReg) increment() {
	r.value += 1
}
