package message_queue

import "sync"

// 每一次触发都会使指针向后面移动
type MessageRing struct {
	Now        []Msg `用来存放没有成功推送出去的消息,消息都会立即推送出去，没推出去下方在Dead`
	Dead       []*Msg
	All        int `当前所有的消息`
	sync.Mutex `互斥锁`
}
type Ring struct {
	Rings [4096]*MessageRing
	N     int
}
type Msg struct {
	From    string
	To      string
	Time    string
	Content string
	Repeat  int
	Loop    int
}

// TODO Message
// 要 new 出一个环形的msg
func NewMessage() *Ring {
	return &Ring{
		Rings: make([]*MessageRing, 1000),
	}
}
