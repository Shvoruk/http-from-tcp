package main

import (
	"fmt"
	"http-from-tcp/internal/request"
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

	fmt.Println("Listening on port 42069")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to establish connection: %v\n", err)
			continue
		}
		fmt.Println("Connection accepted")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("Failed to parse request: %v\n", err)
			conn.Close()
			continue
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
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
