package ppu

import (
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

type SramReader interface {
	ReadSram() [64]objectAttributeEntity
	ReadNextLine() [8]objectAttributeEntity
	ReadCurrentLine() [8]objectAttributeEntity
}

type SramReaderServer struct {
	reader SramReader
}

func NewSramReaderServer(reader SramReader) *SramReaderServer {
	return &SramReaderServer{reader: reader}
}

func (s *SramReaderServer) Handler(ws *websocket.Conn) {
	fmt.Println("Sram reader server: Connection with client: ", ws.RemoteAddr())
	s.connectionHandler(ws)
}

func (s *SramReaderServer) connectionHandler(ws *websocket.Conn) {
	now := time.Now()
	for {
		elapsed := time.Now()
		if elapsed.Sub(now) > FRAME_DURATION {
			oam := s.reader.ReadSram()
			for _, oae := range oam {
				ws.Write([]byte(oae.ToString()))
			}
			now = elapsed
		}
	}
}
