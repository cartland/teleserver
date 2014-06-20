package lib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/calsol/teleserver/can"
	"github.com/gorilla/mux"
)

func frameFromForm(dataType string, r *http.Request) (*can.Frame, error) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return nil, fmt.Errorf("cannot parse id: %v", err)
	}

	switch dataType {
	case "bytes":
		length, err := strconv.Atoi(r.FormValue("length"))
		if err != nil {
			return nil, fmt.Errorf("cannot parse length: %v", err)
		}
		var b []byte
		for i := 0; i < 8 && i < length; i++ {
			formByte, err := strconv.Atoi(r.FormValue(fmt.Sprintf("byte%d", i)))
			if err != nil {
				return nil, fmt.Errorf("cannot parse byte %d: %v", i, err)
			}
			b = append(b, byte(formByte))
		}
		return can.NewFrame(uint32(id), b), nil

	case "floats":
		var floats []float32
		for i := 0; i < 2; i++ {
			float, err := strconv.ParseFloat(r.FormValue(fmt.Sprintf("float%d", i)), 32)
			if err != nil {
				return nil, fmt.Errorf("cannot parse float %d: %v", i, err)
			}
			floats = append(floats, float32(float))
		}
		var buf bytes.Buffer
		for _, f := range floats {
			if err := binary.Write(&buf, binary.LittleEndian, f); err != nil {
				return nil, fmt.Errorf("binary.Write failed: %v", err)
			}
		}
		return can.NewFrame(uint32(id), buf.Bytes()), nil

	default:
		return nil, fmt.Errorf("%s is an invalid send type", dataType)
	}
}

// HandleSendToCAN generates a http handler that sends valid requests over the
// SocketCAN connection.
func HandleSendToCAN(c can.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dataType := mux.Vars(r)["type"]
		frame, err := frameFromForm(dataType, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Got CAN request, want to send %v", frame)
		if c == nil {
			http.Error(w, fmt.Sprintf("CANSocket not running, cannot send %#v over CAN", frame), http.StatusInternalServerError)
			return
		}

		if err := c.WriteFrame(frame); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Sent %v", frame)
		fmt.Fprintf(w, "Successfully wrote data %#v", frame)
	}
}
