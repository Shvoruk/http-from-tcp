package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection accepted")
		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Printf("%s\n", line)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	var ch = make(chan string)

	go func() {
		defer close(ch)
		defer fmt.Printf("Connection closed\n")

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
