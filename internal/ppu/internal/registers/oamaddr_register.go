package registers

type OamAddressReg struct {
	value uint16
}

func (r *OamAddressReg) Reset() {
	r.value = 0
}

func (r *OamAddressReg) Write(value byte) {
	r.value = uint16(value) << 8
}

func (r *OamAddressReg) Increment() {
	r.value += 1
}

func (r *OamAddressReg) GetAddress() uint16 {
	return r.value
}
