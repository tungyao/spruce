package test

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
)

func init() {
	flag.BoolVar(&keep, "keep", false, "keep alive connect")
	flag.StringVar(&addr, "addr", "127.0.0.1:81", "listen port")
}
func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func() {
			c, err := net.Dial("tcp", addr)
			if err != nil {
				log.Panicln(err)
			}
			_, _ = c.Write(EntrySet("a", "a", 0))
			//_, _ = c.Write(EntryGet(key))
			data := make([]byte, 1024)
			_, err = c.Read(data)
			if err != nil {
				log.Println(err)
			}
			//fmt.Println(string(data[:n]))
			_ = c.Close()
			wg.Add(-1)
		}()
	}
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func() {
			c, err := net.Dial("tcp", addr)
			if err != nil {
				log.Panicln(err)
			}
			//_, _ = c.Write(EntrySet("a", "a", 0))
			_, _ = c.Write(EntryGet("a"))
			data := make([]byte, 1024)
			_, err = c.Read(data)
			if err != nil {
				log.Println(err)
			}
			//fmt.Println(string(data[:n]))
			_ = c.Close()
			wg.Add(-1)
		}()
	}
	wg.Wait()

	//tcpAddr, err := net.ResolveTCPAddr("tcp", addr) //创建 tcpAddr数据
	//if err != nil {
	//	log.Println(err,1)
	//}
	//a, err := net.ListenTCP("tcp", tcpAddr)
	//if err != nil {
	//	log.Println(err,2)
	//	return
	//}

	//for {
	//	var operation string
	//	var key string
	//	var value string
	//	fmt.Print(addr + ">> ")
	//	c, err := net.Dial("tcp", addr)
	//	if err != nil {
	//		log.Panicln(err)
	//	}
	//	_, _ = fmt.Scanln(&operation, &key, &value)
	//	if value != "" {
	//		_, _ = c.Write(EntrySet(key, value, 0))
	//	} else {
	//		_, _ = c.Write(EntryGet(key))
	//	}
	//	data := make([]byte, 1024)
	//	n, err:= c.Read(data)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	fmt.Println(string(data[:n]))
	//	_ = c.Close()
	//}
	//for {
	//listen, err := a.AcceptTCP()
	//listen, err := a.A()
	//if err != nil {
	//	log.Println(err,3)
	//	return
	//}
	//go func(c *net.TCPConn) {
	//	var operation string
	//	var key string
	//	var value string
	//	fmt.Print(addr + ">> ")
	//	_, _ = fmt.Scanln(&operation, &key, &value)
	//	if value != "" {
	//		_, _ = c.Write(EntrySet(key, value, 0))
	//	} else {
	//		_, _ = c.Write(EntryGet(key))
	//	}
	//	data := make([]byte, 1024)
	//	n, _ := c.Read(data)
	//	fmt.Println(string(data[:n]))
	//	_ = c.Close()
	//}(listen)

	//}

}
func ParsingExpirationDate(tm interface{}) interface{} {
	switch tm.(type) {
	case []byte:
		if len(tm.([]byte)) > 2 {
			fmt.Println("input error")
		}
		var out int64 = 0
		out = int64(tm.([]byte)[1])
		out += int64(tm.([]byte)[0]) << 8
		return out
	case int:
		out := make([]byte, 2)
		out[1] = byte(tm.(int))
		out[0] = byte(tm.(int) >> 8)
		return out
	}
	return nil
}
func EntrySet(key, value string, ti int) []byte {
	out := make([]byte, 11)
	out[0] = 1
	out[1] = 2
	tm := ParsingExpirationDate(ti).([]byte)
	out[2] = tm[0]
	out[3] = tm[1]
	for _, v := range key {
		out = append(out, byte(v))
	}
	out = append(out, 0)
	for _, v := range value {
		out = append(out, byte(v))
	}
	return out
}
func EntryGet(key string) []byte {
	out := make([]byte, 11)
	out[0] = 2
	out[1] = 2
	for _, v := range key {
		out = append(out, byte(v))
	}
	return out
}
