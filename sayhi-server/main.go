package main

import (
	"log"
	"errors"
	"net"
	"strings"
	"bufio"
	"github.com/sevaaadev/sayhi/internal/protocol"
)

type User struct{
	Name string
	Conn net.Conn
}

func (db UsersOnline) String() string{
	var sb strings.Builder
	for k, _ := range db {
		sb.WriteString(k + "\n")
	}
	return sb.String()
}

type UsersOnline map[string]*User

type State struct {
	db UsersOnline 
}

func relay(data State, out chan protocol.Message) {
	 for {
		 msg := <-out
		 switch msg.Type {
		 case protocol.LogoutMessage:
			 user := data.db[msg.From]
			 user.Conn.Close()
			 delete(data.db, msg.From)
			log.Printf("disconnected from %s\n", msg.From)
		 case protocol.UserMessage:
			 for k, v := range data.db {
				if k != msg.From {
					protocol.WriteMessage(v.Conn, msg)
				}
			 }
			 log.Printf("broadcasting message from %s", msg.From)
		 }
	 }
}

func Auth(conn net.Conn, data State) (*User, error) {
	scanner := bufio.NewScanner(conn)
	scanner.Split(protocol.ScanMessage)
	scanner.Scan()
	msgStruct, _ := protocol.BytesToMessage(scanner.Bytes())
	if msgStruct.Type != protocol.LoginMessage {
		return nil, errors.New("expecting login message") 

	}
	name := msgStruct.Data
	_, ok := data.db[name]
	if ok {
		return nil, errors.New("username is not available, pick another username") 
	}
	user := User{
		Name: name,
		Conn: conn,
	}
	data.db[name] = &user
	log.Printf("connected to %s\n", user.Name)
	return &user, nil 
}

func handleConn(user *User, out chan protocol.Message) {
	scanner := bufio.NewScanner(user.Conn)
	scanner.Split(protocol.ScanMessage)
	for scanner.Scan() {
		msgStruct, err := protocol.BytesToMessage(scanner.Bytes())
		if err != nil {
			log.Printf("WARNING: %s's message cant be decoded: %s", user.Name, err)
			continue
		}
		log.Printf("%s says %s\n", user.Name, msgStruct.Data)
		relayMsg := protocol.Message{
			Type: protocol.UserMessage,
			From: user.Name,
		        Data: msgStruct.Data,
		}
		out <- relayMsg
	}
	logoutMsg := protocol.Message{
		Type: protocol.LogoutMessage,
		From: user.Name,
	}
	out <- logoutMsg
}

const PORT = "7777"

func main() {
	ln, err := net.Listen("tcp4", ":"+PORT)
	if err != nil {
		log.Fatalf("could not listen on port ':%s': %s\n", PORT, err)
	}
	defer ln.Close()
	log.Printf("listening for connection on port :%s\n", PORT)
	data := State{
		db: UsersOnline{},
	}
	out := make(chan protocol.Message)
	go relay(data, out)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("could not accept a connection: %s\n", err)
			continue
		}
		user, err := Auth(conn, data)
		if err != nil {
			log.Printf("fail authentication: %s", err)
			msg := protocol.Message{
				Type: protocol.ErrorMessage,
				Data: err.Error(),
			}
			protocol.WriteMessage(conn, msg)
			conn.Close()
			continue
		}
		go handleConn(user, out)
	}
}
