package test

import (
	"fmt"
	"net"
	"testing"
)

func TestDIS1(t *testing.T) {
	//spruce.StartSpruceDistributed(spruce.Config{Test: ":80"})
	for i := 0; i < 10000; i++ {
		go func() {
			a, _ := net.Dial("tcp", "127.0.0.1:9102")
			a.Write([]byte("get**hello"))
			data := make([]byte, 5)
			a.Read(data)
			fmt.Println(data)
			a.Close()
		}()
	}
}
