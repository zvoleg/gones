package main

import (
	"github.com/zvoleg/gones/internal/bus"
	"github.com/zvoleg/gones/internal/cartridge"
	"github.com/zvoleg/gones/internal/cpu6502"
)

func main() {
	cpuInterruptLine := make(chan cpu6502.Interrupt, 3)
	cartridge := cartridge.New("./gal.nes")
	bus := bus.New(&cartridge)
	cpu := cpu6502.New(&bus, cpuInterruptLine)
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}
}
