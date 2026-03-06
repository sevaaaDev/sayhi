package protocol

import (
	"encoding/json"
	"io"
	"encoding/binary"
)

type MessageType int
const (
	UserMessage MessageType = iota
	ErrorMessage
	LogoutMessage
	LoginMessage
)

type Message struct {
	Type MessageType 
	From string
	Data string
}

func MessageToBytes(msg Message) ([]byte, error) {
	b, err := json.Marshal(msg)
	return b, err
}

func BytesToMessage(msgBytes []byte) (Message, error) {
	var msg Message
	err := json.Unmarshal(msgBytes, &msg)
	return msg, err
}

func WriteMessage(w io.Writer, msg Message) error {
	b, err := MessageToBytes(msg)
	if err != nil {
		return err
	}
	binary.Write(w, binary.BigEndian, uint16(len(b)))
	w.Write(b)
	return nil
}
