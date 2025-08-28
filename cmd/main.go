package main

import (
	"net/http"
	"sync"

	"github.com/zvoleg/gones/internal/bus"
	"github.com/zvoleg/gones/internal/cartridge"
	"github.com/zvoleg/gones/internal/controller"
	"github.com/zvoleg/gones/internal/cpu6502"
	"github.com/zvoleg/gones/internal/ppu"
	"golang.org/x/net/websocket"
)

func main() {
	cartridge := cartridge.New("./smb.nes")
	ppuEmu := ppu.NewPpu()
	joypad := controller.NewJoypad()
	bus := bus.New(&cartridge, &ppuEmu, &joypad)
	ppuEmu.InitBus(&bus)
	cpu := cpu6502.New(&bus)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		server := ppu.NewGuiServer(&ppuEmu)
		http.Handle("/frame", websocket.Handler(server.Handler))
		http.Handle("/pallette", websocket.Handler(server.Handler))
		http.Handle("/pattern", websocket.Handler(server.Handler))
		http.Handle("/name", websocket.Handler(server.Handler))
		http.ListenAndServe(":3000", nil)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		http.HandleFunc("/", render_page)
		http.ListenAndServe(":8081", nil)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		server := controller.NewControllerServer(&joypad)
		http.Handle("/input", websocket.Handler(server.Handler))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clock_counter := 0
		for {
			ppuEmu.Clock()
			if clock_counter%2 == 0 {
				cpu.Clock()
			}
		}
	}()

	wg.Wait()
}

func render_page(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}
