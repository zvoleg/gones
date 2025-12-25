package registers

type Flag byte

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

func (r *StatusReg) SetStatusFlag(f Flag, set bool) {
	if set {
		r.value |= byte(f)
	} else {
		r.value &= ^byte(f)
	}
}
