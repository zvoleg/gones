package ppu

type maskReg struct {
	value byte
}

func (r *maskReg) write(value byte) {
	r.value = value
}

func (r *maskReg) grayscaleDisplayEnabled() bool {
	return (r.value & 0x1) == 1
}

func (r *maskReg) leftBackgroundEnabled() bool {
	return (r.value & 0x2) != 0
}

func (r *maskReg) leftSpritesEnabled() bool {
	return (r.value & 0x4) != 0
}

func (r *maskReg) backgroundEnabled() bool {
	return (r.value & 0x8) != 0
}

func (r *maskReg) spritesEnabled() bool {
	return (r.value & 0x10) != 0
}

// 0x20 Emphasize red (green on PAL/Dendy)
// 0x40 Emphasize green (red on PAL/Dendy)
// 0x80 Emphasize blue
