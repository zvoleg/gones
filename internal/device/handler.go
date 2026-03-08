package device

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Comand struct {
	Comand string `json:"comand"`
}

func (d *Device) HandlePlayerComands(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Panicln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var comand Comand
	err = json.Unmarshal(bodyBytes, &comand)

	if err != nil {
		log.Panicln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch comand.Comand {
	case "playPause":
		d.isRun = !d.isRun
	case "reset":
		d.cpu.Reset()
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{}")
}
