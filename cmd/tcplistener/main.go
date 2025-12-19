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
		fmt.Println("--CONNECTION ESTABLISHED--")

		receivedLines := getLinesChannel(conn)
		for line := range receivedLines {
			fmt.Printf("%s\n", line)
		}

		fmt.Println("--CONNECTION CLOSED--")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {

	ch := make(chan string)
	currentLine := ""

	go func() {
		for {
			data := make([]byte, 8)
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

			if count == 8 {
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
