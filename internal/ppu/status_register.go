package ppu

type flag byte

const (
	O = 0x20
	S = 0x40
	V = 0x80
)

type statusReg struct {
	value byte
	latch *bool
}

func (r *statusReg) read() byte {
	*r.latch = false
	return r.value
}

func (r *statusReg) setStatusFlag(f flag, set bool) {
	if set {
		r.value |= byte(f)
	} else {
		r.value &= ^byte(f)
	}
}
