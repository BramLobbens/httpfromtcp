package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}

receivedLines := getLinesChannel(file)
	for line := range receivedLines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {

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

}
