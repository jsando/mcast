package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"flag"
	"golang.org/x/net/ipv4"
)

func main() {
	hex := flag.Bool("x", false, "Output hexdump -C style output")
	raw := flag.Bool("r", false, "Output raw bytes received, suitable to pipe to a file or other program")
	sendPing := flag.Bool("p", false, "Send ping packets with counter to the same group via the specified adapter (en0, eth0, etc)")
	flag.Usage = func() {
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

func ping(address string) {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	packetConn := ipv4.NewPacketConn(conn)
	mcastLoop, err := packetConn.MulticastLoopback()
	if err != nil {
		panic(err)
	}
	err = packetConn.SetMulticastTTL(32)
	if err != nil {
		log.Fatal("Failed to set multicaset ttl: ", err)
	}
	log.Printf("Multicast loopback: %v", mcastLoop)
	mcastTtl, err := packetConn.MulticastTTL()
	if err != nil {
		panic(err)
	}
	log.Printf("Multicast TTL: %d", mcastTtl)
	dst, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalf("Failed to resolve udp address: %s", address)
	}
	count := uint64(0)
	for {
		msg := fmt.Sprintf("%6d: ping\n", count)
		packetConn.WriteTo([]byte(msg), nil, dst)
		log.Printf("Sent: %s", msg)
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
