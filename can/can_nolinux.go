// +build !linux

package can

import "log"

// Dial requires SocketCAN, a Linux kernel module. It will panic on a
// non-linux system.
func Dial(ifname string) (Conn, error) {
	log.Fatal("Quitting, SocketCAN is not implemented for this platform")
	return nil, nil
}
