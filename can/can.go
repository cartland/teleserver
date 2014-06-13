package can

// +build !linux

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
const (
	CAN_RAW   = iota + 1 /* RAW sockets */
	CAN_BCM              /* Broadcast Manager */
	CAN_TP16             /* VAG Transport Protocol v1.6 */
	CAN_TP20             /* VAG Transport Protocol v2.0 */
	CAN_MCNET            /* Bosch MCNet */
	CAN_ISOTP            /* ISO 15765-2 Transport Protocol */
	CAN_NPROTO
)

// Frame represents a single CAN frame.
// This matches struct can_frame in include/uapi/linux/can.h
type Frame struct {
	Id      uint32 // The CAN ID of the frame and CAN_*_FLAG flags
	DataLen uint8  // Data length code: 0 .. 8
	padding [3]byte
	Data    [8]uint8
}
