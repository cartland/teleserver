package can

import "reflect"

const (
	// special address description flags for the CAN_ID
	CAN_EFF_FLAG = 0x80000000 // EFF/SFF is set in the MSB
	CAN_RTR_FLAG = 0x40000000 // remote transmission request
	CAN_ERR_FLAG = 0x20000000 // error frame

	// valid bits in CAN ID for frame formats
	CAN_SFF_MASK = 0x000007FF // standard frame format (SFF)
	CAN_EFF_MASK = 0x1FFFFFFF // extended frame format (EFF)
	CAN_ERR_MASK = 0x1FFFFFFF // omit EFF, RTR, ERR flags
)

// FrameSize is the size of the Frame struct.
var FrameSize = reflect.TypeOf(Frame{}).Size()

// Frame represents a single CAN frame.
// This matches struct can_frame in include/uapi/linux/can.h
type Frame struct {
	ID      uint32  // The CAN ID of the frame and CAN_*_FLAG flags
	DataLen uint8   // Data length code: 0 .. 8
	_       [3]byte // Spacing to make sure that everything lines up
	Data    [8]byte
}

func (f Frame) CANID() int      { return int(f.ID) }
func (f Frame) CANData() []byte { return f.Data[:] }

// NewFrame returns a new CAN Frame. Data is truncated after 8 bytes.
func NewFrame(msg Message) *Frame {
	id, data := uint32(msg.CANID()), msg.CANData()
	if len(data) > 8 {
		data = data[:8]
	}
	f := &Frame{ID: id, DataLen: uint8(len(data))}
	copy(f.Data[:], data)
	return f
}

// Message represents a single CAN message with ID and Data.
type Message interface {
	CANID() int
	CANData() []byte
}

// Simple is a simple implementation of a Message.
type Simple struct {
	ID   int
	Data []byte
}

func (s Simple) CANID() int      { return s.ID }
func (s Simple) CANData() []byte { return s.Data }

// Conn holds a connection to the CAN socket.
type Conn interface {
	// ReadFrame reads an entire CAN frame at once and returns it. It's recommended
	// to use this instead of Read.
	ReadFrame() (*Frame, error)
	// WriteFrame writes an entire CAN frame at once. It's recommended to use this
	// instead of Write.
	WriteFrame(f *Frame) error
	// Read reads bytes from the CAN socket. When reading, you should pass in a
	// slice at least 16 bytes long to fit an entire frame.
	Read(b []byte) (int, error)
	// Write writes bytes to the CAN socket. When writing, an entire CAN frame should
	// be passed in at a time.
	Write(b []byte) (int, error)
	// Close closes the underlying socket for the connection
	Close() error
}
