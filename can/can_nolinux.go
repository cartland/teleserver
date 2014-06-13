// +build !linux

package can

import "io"

// Dial requires SocketCAN, a Linux kernel module.
func Dial(ifname string) (io.ReadWriteCloser, error) {
	panic("SocketCAN is implemented for this platform")
}
