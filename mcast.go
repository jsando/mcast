package main

import (
	"encoding/hex"
	"fmt"
	"github.com/docopt/docopt-go"
	"log"
	"net"
	"os"
)

func main() {
	usage := `Multicast Dump.

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
`
	args, _ := docopt.Parse(usage, nil, true, "mcast 1.0", false)
	//fmt.Println(args)
	logpackets := true
	raw := args["-r"].(bool)
	hex := args["-x"].(bool)
	if raw {
		logpackets = false
	}
	address := args["ADDRESS:PORT"].(string)
	listen(address, logpackets, hex, raw)
}

func listen(a string, logpackets bool, hexdump bool, raw bool) {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buff := make([]byte, 65536)
	for {
		count, src, err := conn.ReadFromUDP(buff)
		if err != nil {
			fmt.Println("read fail: ", err)
			os.Exit(1)
		}

		if logpackets {
			log.Println("--> ", src, " (", count, " bytes)")
		}
		if hexdump {
			fmt.Println(hex.Dump(buff[:count]))
		}
		if raw {
			os.Stdout.Write(buff[:count])
		}
	}
}
