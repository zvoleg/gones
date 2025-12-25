package registers

type MaskReg struct {
	value byte
}

func (r *MaskReg) Write(value byte) {
	r.value = value
}

func (r *MaskReg) GrayscaleDisplayEnabled() bool {
	return (r.value & 0x1) == 1
}

func (r *MaskReg) LeftBackgroundEnabled() bool {
	return (r.value & 0x2) != 0
}

func (r *MaskReg) LeftSpritesEnabled() bool {
	return (r.value & 0x4) != 0
}

func (r *MaskReg) BackgroundEnabled() bool {
	return (r.value & 0x8) != 0
}

func (r *MaskReg) SpritesEnabled() bool {
	return (r.value & 0x10) != 0
}

func (r *MaskReg) RenderingEnabled() bool {
	return r.BackgroundEnabled() || r.SpritesEnabled()
}

// 0x20 Emphasize red (green on PAL/Dendy)
// 0x40 Emphasize green (red on PAL/Dendy)
// 0x80 Emphasize blue
