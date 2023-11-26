package ppu

type PpuBus interface {
	PpuRead(address uint16) byte
	PpuWrite(address uint16, data byte)
}

type Ppu struct {
	ControllReg   *ControllReg
	MaskReg       *MaskReg
	StatusReg     *StatusReg
	OamAddressReg *OamAddressReg
	OamDataReg    *OamDataReg

	bus PpuBus
}

func NewPpu() Ppu {
	controllReg := ControllReg{}
	maskReg := MaskReg{}
	statusReg := StatusReg{}
	oamAddressReg := OamAddressReg{}
	oamDataReg := OamDataReg{}
	statusReg.value = 0x80 // TODO remove this, it is stub for the testing cpu purpose, there is set up ppu in an always V blank mode

	return Ppu{
		ControllReg:   &controllReg,
		MaskReg:       &maskReg,
		StatusReg:     &statusReg,
		OamAddressReg: &oamAddressReg,
		OamDataReg:    &oamDataReg,
		bus:           nil,
	}
}

func (ppu *Ppu) InitBus(bus PpuBus) {
	ppu.bus = bus
}

func (ppu *Ppu) Clock() {

}
