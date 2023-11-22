package cpu6502

type flag uint8

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
