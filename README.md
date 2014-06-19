[![Build Status](https://travis-ci.org/CalSol/teleserver.svg?branch=master)](https://travis-ci.org/CalSol/teleserver)

# TeleServer

This is a prototype of a telemetry server for a solar car, capable of displaying
a simplified dashboard on a phone or a more complete dashboard on a laptop or
tablet.

At minimum you must supply `-fake`, `-serial`, or `-can_addr` as a data source.

## Installation
1. Install the most recent version of go from http://golang.org/doc/install
2. Create your $GOPATH with `mkdir $HOME/go && export GOPATH=$HOME/go` (http://golang.org/doc/code.html)
3. Run `go get github.com/calsol/teleserver` to fetch the binary
4. Run the binary (`$GOPATH/bin/teleserver -serial /dev/tty`)
5. Navigate to http://localhost:8080

After changing any files in the /public directory, be sure to run
 `go-bindata -o embedded/assets.go -pkg embedded public/...`. This will
  update the resources embedded in the the binary.

To install go-bindata, run `go get github.com/jteeuwen/go-bindata/...`

### Tips and tricks
1. Check out /debug.html for a table of all message types as they come in.
2. Similarly, /dump.html will spit out each message recieved from the server.
3. If using SocketCAN, /sendcan.html allows you to send messages one at a time.
4. `go test github.com/calsol/teleserver/...` helps verify code correctness.

### Flags
* `-port`: Port for the webserver. Default is 8080.
* `-serial`: Port for the serial uart connection.
* `-baud`: Baud rate for the serial port. Default is 115200.
* `-can_addr`: Port for SocketCAN.
* `-fake`: Ignore ports, serve fake data.
* `-sqlite`: Create a sqlite3 database at this location. Default creates a
  temporary database that only lasts as long as the server.
* `-alsologto`: In addition to logging general messages to stdout, log them to
  this file. Default is to only log to stdout.
* `-use_embedded`: Serve the files generated by go-bindata instead of directly
  from /public. This allows the binary to be standalone. Default is false.

### API
* /api/latest?canid=1536&canid=1537 will give you the latest values for messages
  with ids 1536 and 1536 (0x600 and 0x601 in hex)
* /api/graphs?canid=1536&field=ArrayVoltage&time=3m will give you a graph in
  flot format for the ArrayVoltage field of messages with a canid of 0x600.
  Using multiple can ids or multiple fields will give the intersection of
  everything that matches.

### Documentation
Documentation for most of the functionality: [![GoDoc](https://godoc.org/github.com/CalSol/teleserver/lib?status.png)](https://godoc.org/github.com/CalSol/teleserver/lib)

CAN Documentation: [![GoDoc](https://godoc.org/github.com/CalSol/teleserver/can?status.png)](https://godoc.org/github.com/CalSol/teleserver/can)

Embedded Files Documentation: [![GoDoc](https://godoc.org/github.com/CalSol/teleserver/embedded?status.png)](https://godoc.org/github.com/CalSol/teleserver/embedded)


