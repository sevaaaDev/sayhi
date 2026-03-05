package main

import (
	"log"
	// "fmt"
	"net"
	"strings"
	"bufio"
	"github.com/sevaaadev/sayhi/internal/scan"
	"encoding/binary"
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
	conn.Write(append(append([]byte{0, 2}, []byte("BB")...), append([]byte{0, 5}, []byte("hallo")...)...))
	scanner := bufio.NewScanner(conn)
	scanner.Split(scan.ScanMessage)
	for scanner.Scan() {
		msg := scanner.Text()
		log.Printf("%s says %s\n", addr, msg)
		if msg == ":list" {
			connListStr := connList.String()
			binary.Write(conn, binary.BigEndian, uint16(len(connListStr)))
			conn.Write([]byte(connListStr))
			continue
		}
		for _, v := range connList {
			if v != conn {
				msg = conn.RemoteAddr().String() + ": "  + msg
				binary.Write(v, binary.BigEndian, uint16(len(msg)))
				v.Write([]byte(msg))
			}
		}

	}
	conn.Close()
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
