package ppu

type OamAddressReg struct {
	value byte
}

func (r *OamAddressReg) Write(value byte) {
	r.value = value
}
