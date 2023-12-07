package ppu

import (
	"fmt"
	"math/rand"

	"golang.org/x/net/websocket"
)

const (
	frame       = "frame"
	patterTable = "pattern"
	nameTable   = "nameTable"
	collor      = "collor"
)

type ImageProducer interface {
	GetMainScreen() []byte
	GetPatternTables() []byte
	GetNameTable() []byte
	GetCollorPallete() []byte
}

type GuiServer struct {
	imageProducer ImageProducer
}

func NewServer(imageProducer ImageProducer) *GuiServer {
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
	case patterTable:
		s.patternTableSender(ws)
	case nameTable:
		s.nameTableSender(ws)
	}
}

func (s *GuiServer) frameSender(ws *websocket.Conn) {
	for {
		imgSize := 256 * 240
		var imgBuf []byte = make([]byte, imgSize*4)
		for i := 0; i < imgSize*4; i += 4 {
			dot := byte(rand.Intn(3) / 2)
			imgBuf[i] = dot * byte(rand.Intn(256))
			imgBuf[i+1] = dot * byte(rand.Intn(256))
			imgBuf[i+2] = dot * byte(rand.Intn(256))
			imgBuf[i+3] = 255
		}
		_, err := ws.Write(imgBuf)
		if err != nil {
			fmt.Println(err)
			ws.Close()
			return
		}
	}
}

func (s *GuiServer) patternTableSender(ws *websocket.Conn) {
	for {
		srcImg := s.imageProducer.GetPatternTables()
		_, err := ws.Write(srcImg)
		if err != nil {
			fmt.Println(err)
			ws.Close()
			return
		}
	}
}

func (s *GuiServer) nameTableSender(ws *websocket.Conn) {
	for {
		srcImg := s.imageProducer.GetNameTable()
		_, err := ws.Write(srcImg)
		if err != nil {
			fmt.Println(err)
			ws.Close()
			return
		}
	}
}
