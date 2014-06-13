// +build !linux

package can

import "io"

// Dial requires SocketCAN, a Linux kernel module. It will panic on a
// non-linux system.
func Dial(ifname string) (io.ReadWriteCloser, error) {
	panic("SocketCAN is implemented for this platform")
}
