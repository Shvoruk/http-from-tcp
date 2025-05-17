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
	ch := getLinesChannel(file)
	for line := range ch {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	var ch = make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		var buff = make([]byte, 8)
		var line string

		for {
			bytes, err := f.Read(buff)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			part := string(buff[:bytes])
			parts := strings.Split(part, "\n")

			for _, part := range parts[:len(parts)-1] {
				ch <- line + part
				line = ""
			}
			line += parts[len(parts)-1]
		}
		if line != "" {
			ch <- line
		}
	}()
	return ch
}
