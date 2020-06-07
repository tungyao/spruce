package test

import (
	"fmt"
	"log"
	"net"
	"strconv"
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
	tcpC, ok := conn.(*net.TCPConn)
	if !ok {
		log.Println(ok)
	}
	//er:=tcpC.SetKeepAlive(true)
	//if er!=nil{
	//    log.Println(er)
	//}
	//tcpC.SetKeepAlivePeriod(time.Second * 30)
	d := make([]byte, 1024)
	for {
		n, err := tcpC.Read(d)
		if err != nil {
			log.Println(err)
		}
		//tcpC.CloseRead()
		fmt.Println(string(d[:n]))
		n, err = tcpC.Write(d[:n])
		log.Println(n)
		if err != nil {
			log.Println(err)
		}
	}
}
func ()  {
	
}
