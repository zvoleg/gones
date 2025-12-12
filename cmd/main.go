package main

import (
	"net/http"
	"sync"

	"github.com/zvoleg/gones/internal/controller"
	"github.com/zvoleg/gones/internal/device"
	"github.com/zvoleg/gones/internal/ppu"
	"golang.org/x/net/websocket"
)

func main() {
	device := device.NewDevice("./smb.nes")
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		server := ppu.NewGuiServer(device.GetImageProducer())
		http.Handle("/frame", websocket.Handler(server.Handler))
		http.Handle("/palette", websocket.Handler(server.Handler))
		http.Handle("/pattern", websocket.Handler(server.Handler))
		http.Handle("/name", websocket.Handler(server.Handler))
		http.ListenAndServe(":3000", nil)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/index.html")
		})
		http.ListenAndServe(":8081", nil)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		server := controller.NewControllerServer(device.GetJoypadConnector())
		http.Handle("/input", websocket.Handler(server.Handler))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			device.Clock()
		}
	}()

	wg.Wait()
}
