package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	t, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer t.Close()

	current_line := ""
	for {
		data := make([]byte, 8)
		_, err := t.Read(data)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		if i := bytes.IndexByte(data, '\n'); i != -1 {
			current_line += string(data[:i])
			data = data[i+1:]
			fmt.Printf("read: %s\n", current_line)
			current_line = ""
		}

		current_line += string(data)

	}
}
