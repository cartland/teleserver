#TeleServer

This is a prototype of a telemetry server for a solar car, capable of displaying
a simplified dashboard on a phone or a more complete dashboard on a laptop or
tablet. This is a proof of concept, and currently generates fake data.

##Installation
* Install the most recent version of go from http://golang.org/doc/install
* Run `go get bitbucket.org/stvnrhodes/teleserver` to fetch the binary
* Run the binary (at `$GOPATH/bin/teleserver`)
* Navigate to http://localhost:8080

###Flags
`-port`: Port for the webserver. Default is 8080
