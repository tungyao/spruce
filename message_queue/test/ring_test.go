package test

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
	"testing"
	"time"

	"../../message_queue"
)

func TestRing(t *testing.T) {
	// r := message_queue.NewMessage()
	message_queue.Run()
}
func TestDialRing(t *testing.T) {
	conn, err := jsonrpc.Dial("tcp", "127.0.0.1:89")
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}
	var msg = message_queue.Msg{
		From:    "tong",
		To:      "tung",
		Time:    time.Now().Unix(),
		Content: "hello ring",
		Repeat:  0,
		Loop:    0,
	}
	var ok string
	err = conn.Call("Ring.Push", msg, &ok) // 乘法运算
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Println(ok)

}
