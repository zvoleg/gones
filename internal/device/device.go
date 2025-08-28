package device

import (
	"github.com/zvoleg/gones/internal/bus"
	"github.com/zvoleg/gones/internal/cartridge"
	"github.com/zvoleg/gones/internal/controller"
	"github.com/zvoleg/gones/internal/cpu6502"
	"github.com/zvoleg/gones/internal/ppu"
)

type Device struct {
	cpu    cpu6502.Cpu6502
	ppu    ppu.Ppu
	joypad controller.Joypad

	clockCounter uint64
}

func NewDevice(programPath string) Device {
	program := cartridge.New(programPath)
	ppuEmu := ppu.NewPpu()
	joypad := controller.NewJoypad()
	bus := bus.New(&program, &ppuEmu, &joypad)
	ppuEmu.InitBus(&bus)
	cpu := cpu6502.New(&bus)
	return Device{
		cpu:    cpu,
		ppu:    ppuEmu,
		joypad: joypad,

		clockCounter: 0,
	}
}

func (d *Device) Clock() {
	d.ppu.Clock()
	if d.clockCounter%2 == 0 {
		d.cpu.Clock()
	}
	d.clockCounter += 1
}

func (d *Device) GetImageProducer() ppu.ImageProducer {
	return &d.ppu
}

func (d *Device) GetJoypadConnector() controller.Connector {
	return &d.joypad
}
