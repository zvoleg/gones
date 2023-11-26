package ppu

type MaskReg struct {
	value byte
}

func (r *MaskReg) Write(value byte) {
	r.value = value
}

func (r *MaskReg) grayscaleDisplayEnabled() bool {
	return (r.value & 0x1) == 1
}

func (r *MaskReg) leftBackgroundEnabled() bool {
	return (r.value & 0x2) != 0
}

func (r *MaskReg) leftSpritesEnabled() bool {
	return (r.value & 0x4) != 0
}

func (r *MaskReg) backgroundEnabled() bool {
	return (r.value & 0x8) != 0
}

func (r *MaskReg) spritesEnabled() bool {
	return (r.value & 0x10) != 0
}

// 0x20 Emphasize red (green on PAL/Dendy)
// 0x40 Emphasize green (red on PAL/Dendy)
// 0x80 Emphasize blue
