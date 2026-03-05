package main

import (
	"log"
	"net"
	"bufio"
	"github.com/sevaaadev/sayhi/internal/scan"
)

var connList []net.Conn

func handleConn(conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Printf("connected to %s\n", addr)
	conn.Write(append(append([]byte{0, 2}, []byte("BB")...), append([]byte{0, 5}, []byte("hallo")...)...))
	scanner := bufio.NewScanner(conn)
	scanner.Split(scan.ScanMessage)
	for scanner.Scan() {
		log.Printf("%s says %s\n", addr, scanner.Text())
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
