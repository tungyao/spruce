package main

import (
	"flag"
	"fmt"

	ap "git.yaop.ink/tungyao/awesome-pool"
)

var (
	keep bool
	addr string
)

func init() {
	flag.BoolVar(&keep, "keep", false, "keep alive connect")
	flag.StringVar(&addr, "addr", "127.0.0.1:6998", "listen port")
}
func main() {
	p, _ := ap.NewPool(1, "127.0.0.1:6998")
	for {
		var operation string
		var key string
		var value string
		fmt.Print(addr + ">> ")
		_, _ = fmt.Scanln(&operation, &key, &value)
		x := p.Get()
		if value != "" {
			x.Write(EntrySet(key, value, 0))
		} else {
			x.Write(EntryGet(key))
		}
		n := x.Read()
		fmt.Println(n)
	}
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
