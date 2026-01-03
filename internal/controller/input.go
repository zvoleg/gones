package controller

import (
	"fmt"

	"golang.org/x/net/websocket"
)

type InputInterface interface {
	SendByte(data byte)
	ReadBit() byte
}

type Connector interface {
	SetConnection(ws *websocket.Conn)
	IsClosed() bool
}

type Joypad struct {
	ws       *websocket.Conn
	closed   bool
	register byte
}

func NewJoypad() Joypad {
	return Joypad{closed: true}
}

func (j *Joypad) SetConnection(ws *websocket.Conn) {
	j.ws = ws
	j.closed = false
}

func (j *Joypad) IsClosed() bool {
	return j.closed
}

func (j *Joypad) SendByte(data byte) {
	if j.ws == nil || data == 0 {
		return
	}
	_, err := j.ws.Write([]byte{data})
	if err != nil {
		fmt.Println(err)
		fmt.Println("InputInterface: Can't write message to client")
		j.ws.Close()
		j.closed = true
		return
	}
	buffer := make([]byte, 1)
	_, err = j.ws.Read(buffer)
	if err != nil {
		fmt.Println("InputInterface: Can't read message from client")
		j.ws.Close()
		j.closed = true
		return
	}
	if len(buffer) != 0 {
		j.register = buffer[0]
	}
}

func (j *Joypad) ReadBit() byte {
	data := j.register & 0x1
	j.register >>= 1
	return data
}
