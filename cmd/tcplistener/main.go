package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/BramLobbens/httpfromtcp/internal/request"
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
		fmt.Fprintln(os.Stderr, "--CONNECTION ESTABLISHED--")

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for name, values := range r.Headers {
			fmt.Printf("- %s: %s\n", name, values)
		}

		fmt.Fprintln(os.Stderr, "--CONNECTION CLOSED--")
	}
}
