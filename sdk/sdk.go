package sdk

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func newClient() net.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:6998")
	if err != nil {
		log.Panicln(err)
	}
	return conn
}
func Set(key, value []byte, ti int) []byte {
	conn := newClient()
	out := make([]byte, 11)
	out[0] = 1
	out[1] = 2
	tm := ParsingExpirationDate(ti).([]byte)
	out[2] = tm[0]
	out[3] = tm[1]
	for _, v := range key {
		out = append(out, v)
	}
	out = append(out, 0)
	for _, v := range value {
		out = append(out, v)
	}
	conn.Write(out)
	var buf bytes.Buffer
	_, err := io.Copy(&buf, conn)
	if err != nil {
		conn.Close()
		log.Panicln(err)
	}
	return buf.Bytes()
}
func Get(key []byte) []byte {
	conn := newClient()
	out := make([]byte, 11)
	out[0] = 2
	out[1] = 2
	for _, v := range key {
		out = append(out, byte(v))
	}
	conn.Write(out)
	var buf bytes.Buffer
	_, err := io.Copy(&buf, conn)
	if err != nil {
		conn.Close()
		log.Panicln(err)
	}
	return buf.Bytes()
}
func Delete(key []byte) []byte {
	conn := newClient()
	out := make([]byte, 11)
	out[0] = 4
	out[1] = 2
	for _, v := range key {
		out = append(out, byte(v))
	}
	conn.Write(out)
	var buf bytes.Buffer
	_, err := io.Copy(&buf, conn)
	if err != nil {
		conn.Close()
		log.Panicln(err)
	}
	return buf.Bytes()
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
