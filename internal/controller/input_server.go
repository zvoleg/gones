package controller

import (
	"fmt"

	"golang.org/x/net/websocket"
)

type ControllerServer struct {
	connector Connector
}

func NewControllerServer(c Connector) *ControllerServer {
	return &ControllerServer{connector: c}
}

func (s *ControllerServer) Handler(ws *websocket.Conn) {
	fmt.Println("ControllerServer: Connection with client: ", ws.RemoteAddr())
	ws.PayloadType = websocket.BinaryFrame
	s.connector.SetConnection(ws)
	for !s.connector.IsClosed() {
	}
}
