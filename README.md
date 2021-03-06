[![Build Status](https://travis-ci.org/CalSol/teleserver.svg?branch=master)](https://travis-ci.org/CalSol/teleserver)

# TeleServer

This is a prototype of a telemetry server for a solar car, capable of displaying
a simplified dashboard on a phone or a more complete dashboard on a laptop or
tablet.

At minimum you'll want `-fake`, `-serial`, or `-can_addr` as a data source.
Setting `-sqlite` to a filename is needed if you want to read the data later.

## Installation from source (recommended)
1. Install the most recent version of go from http://golang.org/doc/install.
2. Make sure you have [git](http://git-scm.com/). You'll also need gcc (installed by default in OSX and many Linux distros, tested on Windows with [TDM-GCC](http://tdm-gcc.tdragon.net/))
3. In a terminal or command prompt, create your [$GOPATH](http://golang.org/cmd/go/#hdr-GOPATH_environment_variable) with `export GOPATH=$HOME/gocode`. For Windows, run `set GOPATH=%USERPROFILE%\gocode`. It's a good idea to add this as a permanent environmental variable (Google this to figure out how).
4. Run `go get -u -v github.com/calsol/teleserver` to fetch all the code and compile it into a binary. Both the binary and the source code are installed at GOPATH. You'll need to run this command whenever you want to pull updates for the program.
5. Run the binary (`$GOPATH/bin/teleserver -fake` on Linux/OSX, `%GOPATH%/bin/teleserver -fake` on Windows). For getting data from the car, use the `-serial` flag with the appropriate port instead of the `-fake` flag.
6. Navigate to [http://localhost:8080](http://localhost:8080).

### Flags
* `-port`: Port for the webserver. Default is 8080.
* `-serial`: Port for the serial uart connection. Most often something like `/dev/tty.*` for Linux/OSX and something like `COM12` for Windows.
* `-baud`: Baud rate for the serial port. Default is 115200.
* `-canusb`: Treat the serial port as a [CANUSB dongle](http://www.can232.com/?page_id=16)
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
  flot format for the ArrayVoltage field of messages with a canid of 0x600 in
  the last 3 minutes. Using multiple can ids or multiple fields will give the
  intersection of everything that matches.

### Documentation
* Documentation for most of the functionality: [![GoDoc](https://godoc.org/github.com/CalSol/teleserver/lib?status.png)](https://godoc.org/github.com/CalSol/teleserver/lib)

* CAN Documentation: [![GoDoc](https://godoc.org/github.com/CalSol/teleserver/can?status.png)](https://godoc.org/github.com/CalSol/teleserver/can)

* Embedded Files Documentation: [![GoDoc](https://godoc.org/github.com/CalSol/teleserver/embedded?status.png)](https://godoc.org/github.com/CalSol/teleserver/embedded)

After changing any files in the /public directory, be sure to run
 `go-bindata -o embedded/assets.go -ignore \\.bower\.json -ignore bower_components/marked -ignore \demos -ignore \core-tests -ignore bower_components/highlightjs -nomemcopy -pkg embedded public/...`.
  This will update the resources embedded in the the binary.

To install go-bindata, run `go get github.com/jteeuwen/go-bindata/...`



