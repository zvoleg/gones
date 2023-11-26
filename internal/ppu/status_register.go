package ppu

type flag byte

const (
	O = 0x20
	S = 0x40
	V = 0x80
)

type StatusReg struct {
	value byte
}

func (r *StatusReg) Read() byte {
	return r.value
}

func (r *StatusReg) setStatusFlag(f flag, set bool) {

}
