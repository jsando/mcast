package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"flag"
)

func main() {
	hex := flag.Bool("x", false, "Output hexdump -C style output")
	raw := flag.Bool("r", false, "Output raw bytes received, suitable to pipe to a file or other program")
	sendPing := flag.Bool("p", false, "Send ping packets with counter to the same group")
	flag.Usage = func () {
		fmt.Printf("Usage: mcast [-x | -r] [-p] ADDRESS:PORT\n")
	}
	flag.Parse()

	address := flag.Arg(0)
	fmt.Printf("Address: %s\n", address)

	//fmt.Println(args)
	logpackets := true
	if *raw {
		logpackets = false
	}
	if *sendPing {
		go ping(address)
	}
	listen(address, logpackets, *hex, *raw)
}

func ping(a string) {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	count := uint64(0)
	for {
		msg := fmt.Sprintf("%6d: ping\n", count)
		c.Write([]byte(msg))
		time.Sleep(1 * time.Second)
		count++
	}
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
