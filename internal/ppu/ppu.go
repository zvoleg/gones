package ppu

type PpuBus interface {
	PpuRead(address uint16) byte
	PpuWrite(address uint16, data byte)
}

type Ppu struct {
	bus PpuBus
}

func (ppu *Ppu) Clock() {

}
