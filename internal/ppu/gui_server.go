package ppu

import (
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

const (
	frame       = "frame"
	palette     = "palette"
	patterTable = "pattern"
	nameTable   = "nameTable"
)

const FRAME_DURATION = time.Millisecond * 16 // 1 / 60 sec to nanosec

type ImageProducer interface {
	GetMainScreen() []byte
	GetPatternTables() []byte
	GetNameTable() []byte
	GetColorPalette() []byte
}

type GuiServer struct {
	imageProducer ImageProducer
}

func NewGuiServer(imageProducer ImageProducer) *GuiServer {
	return &GuiServer{imageProducer: imageProducer}
}

func (s *GuiServer) Handler(ws *websocket.Conn) {
	fmt.Println("Gui server: Connection with client: ", ws.LocalAddr())
	ticker := time.NewTicker(FRAME_DURATION)

	defer func() {
		fmt.Println("Gui server: Connection closed for client: ", ws.LocalAddr())
		ws.Close()
	}()
	defer ticker.Stop()

	s.connectionHandler(ws, ticker)
}

func (s *GuiServer) connectionHandler(ws *websocket.Conn, ticker *time.Ticker) {
	buffer := make([]byte, 64)
	n, err := ws.Read(buffer)
	if err != nil {
		fmt.Println("Can't read message from client")
	}
	ws.PayloadType = websocket.BinaryFrame
	guiPart := string(buffer[:n])
	switch guiPart {
	case frame:
		writeToSocket(ws, ticker, s.imageProducer.GetMainScreen)
	case palette:
		writeToSocket(ws, ticker, s.imageProducer.GetColorPalette)
	case patterTable:
		writeToSocket(ws, ticker, s.imageProducer.GetPatternTables)
	case nameTable:
		writeToSocket(ws, ticker, s.imageProducer.GetNameTable)
	}
}

func writeToSocket(ws *websocket.Conn, ticker *time.Ticker, imgGeter func() []byte) {
	for {
		select {
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(time.Second * 1))
			buffer := imgGeter()
			_, err := ws.Write(buffer)
			if err != nil {
				fmt.Println(err)
				return
			}
		case <-ws.Request().Context().Done():
			return
		}
	}
}
