package main

import (
	"net/http"
	"sync"

	"github.com/zvoleg/gones/internal/bus"
	"github.com/zvoleg/gones/internal/cartridge"
	"github.com/zvoleg/gones/internal/cpu6502"
	"github.com/zvoleg/gones/internal/ppu"
	"golang.org/x/net/websocket"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		server := ppu.NewServer()
		http.Handle("/nes", websocket.Handler(server.Handler))
		http.ListenAndServe(":3000", nil)
	}()

	cpuInterruptLine := make(chan cpu6502.Interrupt, 3)
	cartridge := cartridge.New("./gal.nes")
	ppu := ppu.NewPpu(cpuInterruptLine)
	bus := bus.New(&cartridge, &ppu)
	ppu.InitBus(&bus)
	cpu := cpu6502.New(&bus, cpuInterruptLine)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			cpu.Clock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			ppu.Clock()
		}
	}()

	wg.Wait()
}
