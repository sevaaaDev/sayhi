package protocol

import (
	"encoding/binary"
)

func ScanMessage(data []byte, atEOF bool) (int, []byte, error) {
	// if data is partial header and its already eof, just assume no token
	if len(data) < 2 && atEOF {
		return len(data), nil, nil
	}
	// if data is partial header, ask for more
	if len(data) < 2 {
		return 0, nil, nil
	}
	msgSize := int(binary.BigEndian.Uint16(data[0:2]))
	// if msg in data is less than msgSize, ask for more data
	if len(data) - 2 < msgSize {
		return 0, nil, nil
	}
	return 2 + msgSize, data[2:msgSize+2:msgSize+2], nil
}
