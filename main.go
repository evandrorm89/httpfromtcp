package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	receiver := getLinesChannel(f)

	for line := range receiver {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	current_line := ""
	go func() {
		defer close(ch)
		for {
			data := make([]byte, 8)
			_, err := f.Read(data)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}

			if i := bytes.IndexByte(data, '\n'); i != -1 {
				current_line += string(data[:i])
				data = data[i+1:]
				ch <- current_line
				current_line = ""
			}

			current_line += string(data)

		}
	}()

	return ch
}
