# mcast
Simple command line utility to show multicast udp traffic.

Requires go (see http://golang.org), with GOPATH setup as per usual.

```
go get github.com/docopt/docopt-go
go install
```

Usage:

```
$ mcast --help
Multicast Dump.

Usage:
  mcast [-x | -r] ADDRESS:PORT
  mcast -h | --help
  mcast --version

Arguments:
  ADDRESS:PORT  The multicast address and port to listen on.

Options:
  -h --help     Show this screen.
  --version     Show version.
  -x            Output hexdump -C style output instead of raw bytes.
  -r            Output raw bytes, suitable to pipe to a file or other program.

If neither -x or -r is specified, it logs the source address of the packets and the number of bytes.
```

There are 3 main usage scenarios:

* Run without "-x" or "-r", it will output the timestamp, origin (remote address), and byte count for each packet.  Useful for logging over time whether packets are arriving consistently.
* Run with "-x" to produce hexdump-style output, to inspect both raw bytes as well as ASCII decoding.
* Run with "-r" to dump raw bytes to stdout, useful to pipe into another program or to capture raw output to a file for later analysis.

