package main

import (
	"io"
	"encoding/binary"
	"os"
	"net"
	"fmt"
	"github.com/sevaaadev/sayhi/internal/scan"
	"bufio"
)

func handleReading(conn io.Reader) {
	scanner := bufio.NewScanner(conn)
	scanner.Split(scan.ScanMessage)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Printf("LOG: the msg is %s\n", msg)
	}
	fmt.Printf("LOG: the server closes the connection\n")
	os.Exit(2)
}

func main() {
	conn, err := net.Dial("tcp", ":7777")
	if err != nil {
		fmt.Printf("ERROR: could not connect to port :7777 : %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	go handleReading(conn)
	inputScanner := bufio.NewScanner(os.Stdin)
	for inputScanner.Scan() {
		input := inputScanner.Text()
		if input == ":q" {
			break
		}
		binary.Write(conn, binary.BigEndian, uint16(len(input)))
		conn.Write(inputScanner.Bytes())
	}
}
