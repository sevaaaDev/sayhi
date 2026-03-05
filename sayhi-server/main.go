package main

import (
	"log"
	"slices"
	"net"
	"strings"
	"bufio"
	"github.com/sevaaadev/sayhi/internal/protocol"
)

type Conns []net.Conn

func (connList Conns) String() string{
	var sb strings.Builder
	for _, v := range connList {
		sb.WriteString(v.RemoteAddr().String() + "\n")
	}
	return sb.String()
}

var connList Conns

func handleConn(conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Printf("connected to %s\n", addr)
	scanner := bufio.NewScanner(conn)
	scanner.Split(protocol.ScanMessage)
	for scanner.Scan() {
		msgStruct, err := protocol.BytesToMessage(scanner.Bytes())
		if err != nil {
			log.Printf("WARNING: %s's message cant be decoded: %s", conn.RemoteAddr(), err)
			continue
		}
		log.Printf("%s says %s\n", addr, msgStruct.Message)
		if msgStruct.Message == ":list" {
			response := protocol.Message{
				From: "server",
				Message: connList.String(),
			}
			protocol.WriteMessage(conn, response)
			continue
		}
		for _, v := range connList {
			if v != conn {
				relayMsg := protocol.Message{
					From: conn.RemoteAddr().String(),
					Message: msgStruct.Message,
				}
				protocol.WriteMessage(v, relayMsg)
			}
		}

	}
	conn.Close()
	connList = slices.DeleteFunc(connList, func(c net.Conn) bool {
		if c == conn {
			return true
		}
		return false
	})
	log.Printf("disconnected from %s\n", addr.String())
}

const PORT = "7777"

func main() {
	ln, err := net.Listen("tcp4", ":"+PORT)
	if err != nil {
		log.Fatalf("could not listen on port ':%s': %s\n", PORT, err)
	}
	log.Printf("listening for connection on port :%s\n", PORT)
	connList = []net.Conn{}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("could not accept a connection: %s\n", err)
			continue
		}
		connList = append(connList, conn)
		go handleConn(conn)
	}
}
