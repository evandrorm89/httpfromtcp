package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("Listening now on localhost:42069")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("A TCP connection has been accepted")
		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}
	}
}

func getLinesChannel(conn net.Conn) <-chan string {
	ch := make(chan string, 1)

	currentLine := ""
	go func() {
		defer close(ch)
		defer conn.Close()
		for {
			data := make([]byte, 8)
			_, err := conn.Read(data)
			if err != nil {
				break
			}

			if i := bytes.IndexByte(data, '\n'); i != -1 {
				currentLine += string(data[:i])
				data = data[i+1:]
				ch <- currentLine
				currentLine = ""
			}

			currentLine += string(data)

		}
		if len(currentLine) != 0 {
			ch <- currentLine
		}
	}()

	return ch
}
