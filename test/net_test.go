package test

import (
	"fmt"
	"net"
	"testing"
)

func TestNet(t *testing.T) {
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println(conn.LocalAddr())
	_, _ = conn.Write([]byte("HELLO"))
}
