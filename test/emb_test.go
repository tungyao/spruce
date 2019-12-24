package test

import (
	"../../spruce"
	"fmt"
	"log"
	"net"
	"testing"
)

func TestEMB(t *testing.T) {
	d := make([]spruce.DCSConfig, 1)
	d[0] = spruce.DCSConfig{
		Name:     "client0",
		Ip:       "127.0.0.1:88",
		Weigh:    0,
		Password: "",
	}
	c := spruce.StartSpruceDistributed(spruce.Config{
		ConfigType: spruce.MEMORY,
		Addr:       "127.0.0.1:88",
		DCSConfigs: d,
		KeepAlive:  true,
		IsBackup:   false,
		NowIP:      "127.0.0.1:88",
	})
	if c == nil {
		log.Panicln(c)
	}
	a, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Println(err)
	}
	for {
		b, err := a.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			defer b.Close()
			data := make([]byte, 1024)
			n, _ := b.Read(data)
			sp := spruce.SplitString(data[:n], []byte("$"))
			for _, v := range sp {
				fmt.Println(string(v))
			}
			if spruce.Equal(sp[0], []byte("set")) {
				o := c.Set(spruce.EntrySet(sp[1], sp[2], 0))
				b.Write(o)
			}
			if spruce.Equal(sp[0], []byte("get")) {
				o := c.Get(spruce.EntryGet(sp[1]))
				b.Write(o)
			}
		}()

	}
}
