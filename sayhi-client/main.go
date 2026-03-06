package main

import (
	"io"
	"os"
	"net"
	"fmt"
	"github.com/sevaaadev/sayhi/internal/protocol"
	"bufio"
)

func handleReading(conn io.Reader) {
	scanner := bufio.NewScanner(conn)
	scanner.Split(protocol.ScanMessage)
	for scanner.Scan() {
		msg, err := protocol.BytesToMessage(scanner.Bytes())
		if err != nil {
			fmt.Printf("WARNING: could not decode message: %s\n", err)
			continue
		}
		if msg.Type == protocol.ErrorMessage {
			msg.From = "error"
		}
		fmt.Printf("%s: %s\n", msg.From, msg.Data)
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
	loginMsg := protocol.Message{
		Type: protocol.LoginMessage,
		Data: os.Args[1],
	}
	protocol.WriteMessage(conn, loginMsg)
	roomMsg := protocol.Message{
		Type: protocol.JoinRoomMessage,
		Data: os.Args[2],
	}
	protocol.WriteMessage(conn, roomMsg)

	inputScanner := bufio.NewScanner(os.Stdin)
	for inputScanner.Scan() {
		input := inputScanner.Text()
		if input == ":q" {
			break
		}
		msg := protocol.Message{
			Type: protocol.UserMessage,
			From: "",
			Room: os.Args[2],
			Data: input,
		}
		protocol.WriteMessage(conn, msg)
	}
}
