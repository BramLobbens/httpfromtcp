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

	currentLine := ""

	for {
		data := make([]byte, 8)
		count, err := file.Read(data)

		currentData := string(data[:count])
		currentLine += currentData

		if err == io.EOF {
			if len(currentLine) > 0 {
				fmt.Printf("read: %s\n", currentLine)
			}
			os.Exit(0)
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
				fmt.Printf("read: %s", part)
			}
			fmt.Println()
			currentLine = lineParts[len(lineParts)-1]
		}
	}
}
