package registers

type OamDataReg struct {
	value byte
}

func (r *OamDataReg) Read() byte {
	return r.value
}

func (r *OamDataReg) Write(value byte) {
	r.value = value
}
