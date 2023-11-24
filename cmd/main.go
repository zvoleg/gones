package main

import (
	"github.com/zvoleg/gones/internal/bus"
	"github.com/zvoleg/gones/internal/cpu6502"
)

func main() {
	cpuInterruptLine := make(chan cpu6502.Interrupt, 3)
	bus := &bus.Bus{}
	cpu := cpu6502.New(bus, cpuInterruptLine)
	cpu.Clock()
}
