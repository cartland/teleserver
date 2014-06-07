# TeleServer

This is a prototype of a telemetry server for a solar car, capable of displaying
a simplified dashboard on a phone or a more complete dashboard on a laptop or
tablet. This is a proof of concept, and currently generates fake data.

## Installation
1. Install the most recent version of go from http://golang.org/doc/install
2. Run `go get github.org/calsol/teleserver` to fetch the binary
3. Run the binary (`$GOPATH/bin/teleserver -serial /dev/tty`)
4. Navigate to http://localhost:8080

### Flags
* `-port`: Port for the webserver. Default is 8080.
* `-serial`: Port for the serial uart connection. Will error if not provided.
* `-baud`: Baud rate for the serial port. Default is 115200.
* `-can`: Interpret the serial input as binary CAN messages instead of JSON.
   Default is false.
* `-fake`: Ignore serial port, serve fake data.
* `-log_file`: Log to a file with this prefix. "_YYYY-MM-DD_HH:MM:SS" is added
   as a suffix. Default prefix is "_tmp".
