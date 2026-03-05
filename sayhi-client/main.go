package main

import (
	"os"
	"net"
	"fmt"
	"github.com/sevaaadev/sayhi/internal/scan"
	"bufio"
)

func main() {
	conn, err := net.Dial("tcp", ":7777")
	if err != nil {
		fmt.Printf("ERROR: could not connect to port :7777 : %s", err)
		os.Exit(1)
	}
	sc := bufio.NewScanner(conn)
	sc.Split(scan.ScanMessage)
	for sc.Scan() {
		msg := sc.Text()
		fmt.Printf("LOG: the msg is %s\n", msg)
	}
}
