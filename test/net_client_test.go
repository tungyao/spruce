package test

import (
	"log"
	"net"
	"testing"
)

func TestNetClient(t *testing.T) {
	conn, _ := net.Dial("tcp", "localhost:8000")
	defer conn.Close()
	buf := make([]byte, 1024*8)
	_, _ = conn.Read(buf)
	t.Log(string(buf))
}
func TestIDontKnow(t *testing.T) {
	a := []byte{0xff}
	log.Println(a)
}
