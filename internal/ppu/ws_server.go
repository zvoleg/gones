package ppu

import (
	"fmt"
	"math/rand"

	"golang.org/x/net/websocket"
)

type Server struct {
	ws *websocket.Conn
}

func NewServer() *Server {
	return &Server{ws: nil}
}

func (s *Server) Handler(ws *websocket.Conn) {
	fmt.Println("Connection with client: ", ws.RemoteAddr())
	s.ws = ws
	s.ws.PayloadType = websocket.BinaryFrame
	s.writeLoop()
}

func (s *Server) writeLoop() {
	for {
		imgSize := 256 * 244
		var imgBuf []byte = make([]byte, imgSize*4)
		for i := 0; i < imgSize*4; i += 4 {
			dot := byte(rand.Intn(3) / 2)
			imgBuf[i] = dot * byte(rand.Intn(256))
			imgBuf[i+1] = dot * byte(rand.Intn(256))
			imgBuf[i+2] = dot * byte(rand.Intn(256))
			imgBuf[i+3] = 255
		}
		s.ws.Write([]byte(imgBuf))
	}
}
