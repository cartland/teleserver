package lib

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/calsol/teleserver/messages"
	"github.com/stvnrhodes/broadcaster"
)

const (
	startByte      = 0xE7
	escapeByte     = 0x75
	idMask         = 0xFFF0
	lenMask        = 0x000F
	maxPayloadSize = 8
)

// unescape implements decoding data streams where the data is escaped by XORing
// significant bytes with an escape byte.
func unescape(bs []byte) []byte {
	var fixed []byte
	for {
		i := bytes.IndexByte(bs, escapeByte)
		if i < 0 {
			return append(fixed, bs...)
		}
		fixed = append(fixed, bs[:i]...)
		if len(bs) < i+1 {
			// The last byte was an encoded value. Wrong, but we'll roll with it.
			return fixed
		}
		fixed = append(fixed, bs[i]^bs[i+1])
		bs = bs[i+2:]
	}
}

// NewXSPScanner creates a scanner that divides up the input at every 0xE7 byte
// and decodes it at every 0x75 byte.
func NewXSPScanner(r io.Reader) *bufio.Scanner {
	s := bufio.NewScanner(r)
	s.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		if i := bytes.IndexByte(data, startByte); i >= 0 {
			// We have a full newline-terminated line.
			return i + 1, unescape(data[0:i]), nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), unescape(data), nil
		}
		// Request more data.
		return 0, nil, nil
	})
	return s
}

func newCANFromBytes(b []byte) (messages.CAN, error) {
	if len(b) < 2 {
		return nil, fmt.Errorf("message was too short: %v", b)
	}
	header, body := b[:2], b[2:]
	preamble := binary.LittleEndian.Uint16(header)
	id := (preamble & idMask) >> 4
	length := int(preamble & lenMask)

	if length > maxPayloadSize {
		return nil, fmt.Errorf("packet 0x%x: payload size %d is greater than %d: %v", id, length, maxPayloadSize, body)
	} else if len(body) != length {
		return nil, fmt.Errorf("packet 0x%x: payload size %d != actual size %d: %v", id, length, len(body), body)
	}
	msg, ok := messages.IDToMessage[id]
	if !ok {
		return nil, fmt.Errorf("packet 0x%x: unknown id, size %d: %v", id, length, body)
	}

	// Make a new copy of the message to prevent changing the one in the id map.
	msg = msg.New()

	if err := msg.UnmarshalBinary(body); err != nil {
		return nil, fmt.Errorf("packet 0x%x: payload %v: %v", id, body, err)
	}
	return msg, nil
}

// CANReader allows reading CAN messages from an io.Reader.
type CANReader struct {
	b *bufio.Scanner
}

// NewCanReader creates a reader that scans the reader as XSP and parses as CAN.
func NewCANReader(r io.Reader) *CANReader {
	return &CANReader{b: NewXSPScanner(r)}
}

// Read returns the next CAN message from the input.
func (c *CANReader) Read() (messages.CAN, error) {
	if c.b.Scan() {
		return newCANFromBytes(c.b.Bytes())
	}
	return nil, c.b.Err()
}

// ReadCAN will continually read bytes from the io.Reader, interpret them as
// binary CAN messages, and send them through the broadcaster.
func ReadCAN(r io.Reader, b broadcaster.Caster) {
	c := NewCANReader(r)
	for {
		msg, err := c.Read()
		if err != nil {
			log.Print(err)
			continue
		}
		b.Cast(messages.CANPlus{msg, time.Now()})
	}
}
