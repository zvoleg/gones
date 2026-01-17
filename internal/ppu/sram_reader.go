package ppu

func (ppu *Ppu) ReadSram() [64]objectAttributeEntity {
	return ppu.oam
}

func (ppu *Ppu) ReadNextLine() [8]objectAttributeEntity {
	return ppu.nextLineOam
}

func (ppu *Ppu) ReadCurrentLine() [8]objectAttributeEntity {
	return ppu.curentLineOam
}
