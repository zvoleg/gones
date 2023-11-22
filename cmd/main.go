package main

import (
	"github.com/zvoleg/gones/internal/bus"
	"github.com/zvoleg/gones/internal/cpu6502"
)

func main() {
	bus := &bus.Bus{}
	cpu := cpu6502.New(bus)
	cpu.Clock()
}
