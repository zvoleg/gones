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
	collor      = "collor"
)

const FRAME_DURATION = time.Duration(16666666) // 1 / 60 sec to nanosec

type ImageProducer interface {
	GetMainScreen() []byte
	GetPatternTables() []byte
	GetNameTable() []byte
	GetCollorPalette() []byte
}

type GuiServer struct {
	imageProducer ImageProducer
}

func NewGuiServer(imageProducer ImageProducer) *GuiServer {
	return &GuiServer{imageProducer: imageProducer}
}

func (s *GuiServer) Handler(ws *websocket.Conn) {
	fmt.Println("Connection with client: ", ws.RemoteAddr())
	s.connectionHandler(ws)
}

func (s *GuiServer) connectionHandler(ws *websocket.Conn) {
	buffer := make([]byte, 64)
	n, err := ws.Read(buffer)
	if err != nil {
		fmt.Println("Can't read message from client")
	}
	ws.PayloadType = websocket.BinaryFrame
	guiPart := string(buffer[:n])
	switch guiPart {
	case frame:
		s.frameSender(ws)
	case palette:
		s.paletteSender(ws)
	case patterTable:
		s.patternTableSender(ws)
	case nameTable:
		s.nameTableSender(ws)
	}
}

func (s *GuiServer) frameSender(ws *websocket.Conn) {
	imgBuf := s.imageProducer.GetMainScreen()
	now := time.Now()
	for {
		elapsed := time.Now()
		if elapsed.Sub(now) > FRAME_DURATION {
			// TODO rendering by signal from ppu (frame_done) and sync with time 1/60 sec
			_, err := ws.Write(imgBuf)
			if err != nil {
				fmt.Println(err)
				ws.Close()
				return
			}
			now = elapsed
		}
	}
}

func (s *GuiServer) paletteSender(ws *websocket.Conn) {
	now := time.Now()
	for {
		elapsed := time.Now()
		if elapsed.Sub(now) > FRAME_DURATION {
			srcImg := s.imageProducer.GetCollorPalette()
			_, err := ws.Write(srcImg)
			if err != nil {
				fmt.Println(err)
				ws.Close()
				return
			}
			now = elapsed
		}
	}
}

func (s *GuiServer) patternTableSender(ws *websocket.Conn) {
	now := time.Now()
	for {
		elapsed := time.Now()
		if elapsed.Sub(now) > FRAME_DURATION {
			srcImg := s.imageProducer.GetPatternTables()
			_, err := ws.Write(srcImg)
			if err != nil {
				fmt.Println(err)
				ws.Close()
				return
			}
			now = elapsed
		}
	}
}

func (s *GuiServer) nameTableSender(ws *websocket.Conn) {
	now := time.Now()
	for {
		elapsed := time.Now()
		if elapsed.Sub(now) > FRAME_DURATION {
			srcImg := s.imageProducer.GetNameTable()
			_, err := ws.Write(srcImg)
			if err != nil {
				fmt.Println(err)
				ws.Close()
				return
			}
			now = elapsed
		}
	}
}
