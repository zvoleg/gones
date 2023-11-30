package ppu

type oamDataReg struct {
	value byte
}

func (r *oamDataReg) read() byte {
	return r.value
}

func (r *oamDataReg) write(value byte) {
	r.value = value
}
