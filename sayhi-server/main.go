package main

import (
	"log"
	"slices"
	"net"
	"strings"
	"bufio"
	"github.com/sevaaadev/sayhi/internal/protocol"
)

type User struct{
	Name string
	Conn net.Conn
}

func StringUser(u []User) string{
	var sb strings.Builder
	for _, v := range u {
		sb.WriteString(v.Name + "\n")
	}
	return sb.String()
}

var UserOnline []User 

func handleConn(conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Printf("connected to %s\n", addr)
	scanner := bufio.NewScanner(conn)
	scanner.Split(protocol.ScanMessage)
	scanner.Scan() 
	msg, _ := protocol.BytesToMessage(scanner.Bytes())
	user := User {
		Name: msg.Message,
		Conn: conn,
	}
	UserOnline = append(UserOnline, user)
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
				Message: StringUser(UserOnline),
			}
			protocol.WriteMessage(conn, response)
			continue
		}
		for _, v := range UserOnline {
			if v.Conn != conn {
				relayMsg := protocol.Message{
					From: user.Name,
					Message: msgStruct.Message,
				}
				protocol.WriteMessage(v.Conn, relayMsg)
			}
		}

	}
	conn.Close()
	UserOnline = slices.DeleteFunc(UserOnline, func(c User) bool {
		if c.Conn == conn {
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
	UserOnline = []User{}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("could not accept a connection: %s\n", err)
			continue
		}
		go handleConn(conn)
	}
}
