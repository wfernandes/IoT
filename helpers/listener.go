package main

import (
	"net"

	"fmt"

	"github.com/pivotal-golang/localip"
)

func main() {

	ip, _ := localip.LocalIP()
	conn, err := net.ListenPacket("udp", net.JoinHostPort(ip, "8888"))
	if err != nil {
		panic(err)
	}

	readBuffer := make([]byte, 1000)
	fmt.Printf("Listening on %s:8888", ip)
	for {
		readCount, _, err := conn.ReadFrom(readBuffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Read %d bytes", readCount)
		fmt.Println(string(readBuffer))
	}
}
