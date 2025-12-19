package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn := setupUDPConnection("42069")
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	var bear rune = '🐻'
	for {
		fmt.Printf("%q", bear)
		readString, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = conn.Write([]byte(readString))
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func setupUDPConnection(port string) *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:"+port)
	if err != nil {
		log.Fatal(err.Error())
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	return conn
}
