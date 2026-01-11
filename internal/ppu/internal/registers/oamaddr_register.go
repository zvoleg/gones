package registers

type OamAddressReg struct {
	value byte
}

func (r *OamAddressReg) Reset() {
	r.value = 0
}

func (r *OamAddressReg) Write(value byte) {
	r.value = value
}

func (r *OamAddressReg) Increment() {
	r.value += 1
}

func (r *OamAddressReg) GetAddress() byte {
	return r.value
}
