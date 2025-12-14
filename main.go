package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	for {
		data := make([]byte, 8)
		count, err := file.Read(data)
		if err == io.EOF {
			os.Exit(0)
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("read: %s\n", data[:count])
	}
}
