package lib

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/calsol/teleserver/binpack"
	"github.com/calsol/teleserver/can"
	"github.com/calsol/teleserver/msgs"
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

// newCANFromBytes takes the raw bytes of a CAN message and parses it into a
// semantically useful message.
func newCANFromBytes(b []byte) (msgs.CAN, error) {
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

	msg := msgs.IDToMessage(id)
	if err := binpack.Unmarshal(body, msg); err != nil {
		return nil, fmt.Errorf("packet 0x%x: payload %v: %v", id, body, err)
	}
	return msg, nil
}

// CANReader is capable of reading CAN messages.
type CANReader interface {
	// Read gets the next message from the CANReader
	Read() (msgs.CAN, error)
}

// xspCANReader allows reading CAN messages from a XSP Serial CAN connection.
type xspCANReader struct {
	b *bufio.Scanner
}

// NewXSPCANReader creates a reader that scans as XSP and parses as CAN. The XSP
// format reads a stream of bytes and uses a separator character to split
// individual messages apart
func NewXSPCANReader(r io.Reader) CANReader {
	return &xspCANReader{b: NewXSPScanner(r)}
}
func (c *xspCANReader) Read() (msgs.CAN, error) {
	if c.b.Scan() {
		return newCANFromBytes(c.b.Bytes())
	}
	return nil, c.b.Err()
}

type socketCANReader struct {
	r io.Reader
}

// We fudge the numbers a bit here so that this matches XSPCAN messages. We
// remove the padding, combine the length with the id, and truncate the data.
// If we ever deprecate XSPCAN, newCANFromBytes should probably be changed to
// expect binary encodings of can.Frame.
func changeSocketCANEncoding(b []byte) ([]byte, error) {
	id := binary.LittleEndian.Uint32(b[:4])
	body := b[8:]
	length := b[4]
	if int(length) > len(body) {
		return nil, fmt.Errorf("packet 0x%x: payload size %d is greater than %d: %v", id, len(body), length, body)
	}
	ids := make([]byte, 4)
	binary.LittleEndian.PutUint32(ids, id<<4)
	ids[0] |= length
	return append(ids[:2], body[:length]...), nil
}

// NewSocketCANReader creates a reader that reads from a SocketCAN connection
// and returns complete messages. The SocketCAN format reads can.FrameSize bytes
// at a time and interprets them as a complete CAN message.
func NewSocketCANReader(r io.Reader) CANReader { return socketCANReader{r} }
func (s socketCANReader) Read() (msgs.CAN, error) {
	b := make([]byte, can.FrameSize)
	if n, err := s.r.Read(b); err != nil {
		return nil, err
	} else if n != len(b) {
		return nil, fmt.Errorf("got %d bytes, want %d", n, len(b))
	}
	msg, err := changeSocketCANEncoding(b)
	if err != nil {
		return nil, err
	}
	return newCANFromBytes(msg)
}

// ReadCAN will continually read bytes from the io.Reader, interpret them as
// binary CAN messages, and send them through the broadcaster.
func ReadCAN(r CANReader, b broadcaster.Caster) {
	for {
		msg, err := r.Read()
		if err != nil {
			log.Print(err)
			continue
		}
		b.Cast(msgs.NewCANPlus(msg))
	}
}
