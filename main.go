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

	receiver := getLinesChannel(f)

	for line := range receiver {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string, 1)

	current_line := ""
	go func() {
		defer close(ch)
		defer f.Close()
		for {
			data := make([]byte, 8)
			_, err := f.Read(data)
			if err != nil {
				break
			}

			if i := bytes.IndexByte(data, '\n'); i != -1 {
				current_line += string(data[:i])
				data = data[i+1:]
				ch <- current_line
				current_line = ""
			}

			current_line += string(data)

		}
		if len(current_line) != 0 {
			ch <- current_line
		}
	}()

	return ch
}
