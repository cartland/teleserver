package can

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"sync"
	"syscall"
)

const (
	/* special address description flags for the CAN_ID */
	CAN_EFF_FLAG = 0x80000000 /* EFF/SFF is set in the MSB */
	CAN_RTR_FLAG = 0x40000000 /* remote transmission request */
	CAN_ERR_FLAG = 0x20000000 /* error frame */

	/* valid bits in CAN ID for frame formats */
	CAN_SFF_MASK = 0x000007FF /* standard frame format (SFF) */
	CAN_EFF_MASK = 0x1FFFFFFF /* extended frame format (EFF) */
	CAN_ERR_MASK = 0x1FFFFFFF /* omit EFF, RTR, ERR flags */
)

var FrameSize = reflect.TypeOf(Frame{}).Size()

// Frame represents a single CAN frame.
// This matches struct can_frame in include/uapi/linux/can.h
type Frame struct {
	ID      uint32  // The CAN ID of the frame and CAN_*_FLAG flags
	DataLen uint8   // Data length code: 0 .. 8
	Padding [3]byte // Spacing to make sure that everything lines up
	Data    [8]byte
}

// NewFrame returns a new CAN Frame. Data is truncated after 8 bytes.
func NewFrame(id uint32, data []byte) *Frame {
	if len(data) > 8 {
		data = data[:8]
	}
	f := &Frame{ID: id, DataLen: uint8(len(data))}
	copy(f.Data[:], data)
	return f
}

// Conn holds a connection to the CAN socket.
type Conn struct {
	ifname string
	fd     int
	mu     sync.Mutex
	buf    bytes.Buffer
}

// ReadFrame reads an entire CAN frame at once and returns it. It's recommended
// to use this instead of Read.
func (c *Conn) ReadFrame() (*Frame, error) {
	var f Frame
	if err := binary.Read(c, binary.LittleEndian, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// WriteFrame writes an entire CAN frame at once. It's recommended to use this
// instead of Write.
func (c *Conn) WriteFrame(f *Frame) error {
	return binary.Write(c, binary.LittleEndian, f)
}

// When reading, you should pass in a slice at least 16 bytes long to fit an
// entire frame.
func (c *Conn) Read(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.buf.Len() == 0 {
		buf := make([]byte, 16)
		n, err := syscall.Read(c.fd, buf)
		if err != nil {
			return 0, err
		}
		c.buf.Write(buf[:n])
	}
	return c.buf.Read(b)
}

// When writing, an entire CAN frame should be passed in at a time.
func (c *Conn) Write(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return syscall.Write(c.fd, b)
}

func (c *Conn) Close() error {
	return syscall.Close(c.fd)
}
