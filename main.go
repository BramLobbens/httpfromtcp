package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("CONNECTION ESTABLISHED")

		receivedLines := getLinesChannel(conn)
		for line := range receivedLines {
			fmt.Printf("read: %s\n", line)
		}
	}
}

func getLinesChannel(f net.Conn) <-chan string {

	const BYTE_SIZE int = 8
	ch := make(chan string) //output channel
	currentLine := ""

	go func() {
		for {
			data := make([]byte, BYTE_SIZE)
			count, err := f.Read(data)

			currentLine += string(data[:count])

			if err == io.EOF {
				if len(currentLine) > 0 {
					ch <- currentLine
				}
				close(ch)
			}

			if err != nil {
				log.Fatal(err)
			}

			if count == BYTE_SIZE {
				lineParts := strings.Split(currentLine, "\n")
				partsLength := len(lineParts)

				if partsLength == 1 {
					continue
				}

				for _, part := range lineParts[:len(lineParts)-1] {
					ch <- part
				}

				currentLine = lineParts[len(lineParts)-1]
			}
		}
	}()

	return ch
}
