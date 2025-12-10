package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	address, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sliceInput := []byte(input)
		_, err = conn.Write(sliceInput)
		if err != nil {
			log.Fatal(err)
		}
	}
}
