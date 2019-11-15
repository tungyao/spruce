package test

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestKeepAlive(t *testing.T) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ":88") //创建 tcpAddr数据
	a, _ := net.ListenTCP("tcp", tcpAddr)
	for {
		c, _ := a.AcceptTCP()
		go func() {
			c.SetKeepAlive(true)
			data := make([]byte, 1024)
			n, _ := c.Read(data)
			fmt.Println(string(data[:n]))
		}()

	}
}
func TestClientKeep(t *testing.T) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ":88") //创建 tcpAddr数据
	a, _ := net.DialTCP("tcp", nil, tcpAddr)
	for {
		a.Write([]byte("hello"))
		time.Sleep(time.Second * 1)
	}
}
