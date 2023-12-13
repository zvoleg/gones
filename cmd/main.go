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

	cpuInterruptLine := make(chan cpu6502.Signal, 3)
	cartridge := cartridge.New("./dk.nes")
	ppuEmu := ppu.NewPpu(cpuInterruptLine)
	bus := bus.New(&cartridge, &ppuEmu)
	ppuEmu.InitBus(&bus)
	cpu := cpu6502.New(&bus, cpuInterruptLine)

	wg.Add(1)
	go func() {
		defer wg.Done()

		server := ppu.NewServer(&ppuEmu)
		http.Handle("/frame", websocket.Handler(server.Handler))
		http.Handle("/pallette", websocket.Handler(server.Handler))
		http.Handle("/pattern", websocket.Handler(server.Handler))
		http.Handle("/name", websocket.Handler(server.Handler))
		http.ListenAndServe(":3000", nil)
	}()
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
			ppuEmu.Clock()
		}
	}()

	wg.Wait()
}
