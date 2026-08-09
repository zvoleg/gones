package apu

import (
	"bytes"
	"encoding/binary"
	"math"
	"time"

	"golang.org/x/net/websocket"
)

type ApuServer struct {
}

func NewApuServer() *ApuServer {
	return &ApuServer{}
}

func (s *ApuServer) Handler(ws *websocket.Conn) {
	defer ws.Close()

	ws.PayloadType = websocket.BinaryFrame

	sampleRate := 44100
	frequency := 440
	phaseIncrement := (2.0 * math.Pi * float64(frequency)) / float64(sampleRate)
	phase := 0.0

	bufferSize := 2048
	buffer := new(bytes.Buffer)

	interval := time.Duration(float64(bufferSize)/float64(sampleRate)*1000.0) * time.Millisecond
	nextSend := time.Now().Add(interval)

	for {
		select {
		case <-ws.Request().Context().Done():
			return
		default:
			// audio stub just a sin wave
			for range bufferSize {
				value := math.Sin(phase) * 0.5
				binary.Write(buffer, binary.LittleEndian, float32(value))
				phase += phaseIncrement
				if phase >= 2.0*math.Pi {
					phase = 0.0
				}
			}
			ws.SetWriteDeadline(time.Now().Add(time.Second * 1))
			ws.Write(buffer.Bytes())
			buffer.Reset()

			now := time.Now()
			if now.Before(nextSend) {
				time.Sleep(nextSend.Sub(now))
			}
			nextSend = nextSend.Add(interval)
		}
	}
}
