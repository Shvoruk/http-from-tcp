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
	defer file.Close()

	var buff = make([]byte, 8)
	var line string

	for {
		bytes, err := file.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		part := string(buff[:bytes])
		parts := strings.Split(part, "\n")
		for _, part := range parts[:len(parts)-1] {
			fmt.Printf("read: %s\n", line+part)
			line = ""
		}
		line += parts[len(parts)-1]
	}
	fmt.Printf("read: %s\n", line)
}
