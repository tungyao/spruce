package test

import (
	"fmt"
	"strconv"
	"testing"
)
import "../../message_queue"

func TestQueue(t *testing.T) {
	q := message_queue.InitQueue()
	for i := 0; i < 10; i++ {
		q.InQueue(&message_queue.Body{
			Next:    nil,
			Key:     []byte(strconv.Itoa(i)),
			Value:   []byte(strconv.Itoa(i)),
			Operate: nil,
			Check:   true,
		})
	}
	fmt.Println("---------", q.Body.Next)
	//t.Log(string(q.OutQueue().Key))
	//t.Log(string(q.OutQueue().Key))
	//t.Log(string(q.OutQueue().Key))
}
