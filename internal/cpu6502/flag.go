package cpu6502

type flag byte

const (
	C flag = iota
	Z
	I
	D
	B
	U
	V
	N
)
