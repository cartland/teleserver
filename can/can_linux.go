package can

/*
  The following flags should be part of the config for the kernel

CONFIG_CAN=m
CONFIG_CAN_RAW=m
CONFIG_CAN_BCM=m
CONFIG_CAN_VCAN=m
*/

import (
	"bytes"
	"encoding/binary"
	"log"
	"reflect"
	"sync"
	"syscall"
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
const (
	ifnamsiz = 16 /* Size of the name for ifreq */
)

// This matches struct sockaddr_can in include/uapi/linux/can.h
type sockAddr struct {
	Family  uint16
	padding [2]byte
	Ifindex int32
	Addr    [8]byte
}

// This matches struct ifreq in include/uapi/linux/if.h
type ifReq struct {
	Name  [ifnamsiz]byte
	Union [40 - ifnamsiz]byte
}

// Dial returns an object that is used to communicate with a SocketCAN
// connection. This is a naively implemented solution which returns a connection
// that holds a lock when reading or writing. Do not try opening the same
// connection twice or reusing a closed connection, there's nothing to protect
// you.
func Dial(ifname string) (Conn, error) {
	fd, err := syscall.Socket(syscall.AF_CAN, syscall.SOCK_RAW, CAN_RAW)
	if err != nil {
		log.Println("Failed to open socket")
		return nil, err
	}

	ifr := &ifReq{}
	copy(ifr.Name[:], ifname)
	if err := ioctl(fd, syscall.SIOCGIFINDEX, ifr); err != nil {
		log.Println("Failed to ioctl")
		syscall.Close(fd)
		return nil, err
	}

	addr := &sockAddr{Family: syscall.AF_CAN}
	addr.Ifindex = int32(binary.LittleEndian.Uint32(ifr.Union[:]))
	if err := bind(fd, addr); err != nil {
		log.Println("Failed to bind")
		syscall.Close(fd)
		return nil, err
	}

	return &conn{ifname: ifname, fd: fd}, nil
}

func ptrAndSize(n interface{}) (ptr, size uintptr) {
	return reflect.ValueOf(n).Pointer(), reflect.TypeOf(n).Elem().Size()
}

// errno is an int, we want it to be nil if 0.
func toErr(err syscall.Errno) error {
	if err != 0 {
		return err
	}
	return nil
}

// There's no nice wrapper for ioctl, so let's make one.
func ioctl(fd int, typ uintptr, struc interface{}) error {
	ptr, _ := ptrAndSize(struc)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), typ, ptr)
	return toErr(errno)
}

// Boo, we can't use syscall.Bind because the interface has a private member :-(
func bind(fd int, sockAddr interface{}) error {
	ptr, size := ptrAndSize(sockAddr)
	_, _, errno := syscall.Syscall(syscall.SYS_BIND, uintptr(fd), ptr, size)
	return toErr(errno)
}

// conn holds a connection to the CAN socket.
type conn struct {
	ifname    string
	fd        int
	readLock  sync.Mutex
	writeLock sync.Mutex
	buf       bytes.Buffer
}

// ReadFrame reads an entire CAN frame at once and returns it. It's recommended
// to use this instead of Read.
func (c *conn) ReadFrame() (*Frame, error) {
	var f Frame
	if err := binary.Read(c, binary.LittleEndian, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// WriteFrame writes an entire CAN frame at once. It's recommended to use this
// instead of Write.
func (c *conn) WriteFrame(f *Frame) error {
	return binary.Write(c, binary.LittleEndian, f)
}

// Read reads bytes from the CAN socket. When reading, you should pass in a
// slice at least 16 bytes long to fit an entire frame.
func (c *conn) Read(b []byte) (int, error) {
	c.readLock.Lock()
	defer c.readLock.Unlock()
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

// Write writes bytes to the CAN socket. When writing, an entire CAN frame should
// be passed in at a time.
func (c *conn) Write(b []byte) (int, error) {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	return syscall.Write(c.fd, b)
}

// Close closes the underlying socket for the connection
func (c *conn) Close() error {
	return syscall.Close(c.fd)
}
