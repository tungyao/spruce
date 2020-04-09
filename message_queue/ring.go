package message_queue

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"sync"
)

type Ring struct {
	Now      []*Msg
	Dead     []*Msg `用来存放没有成功推送出去的消息,消息都会立即推送出去，没推出去下方在Dead`
	All      int    `当前所有的消息数`
	mux      sync.RWMutex
	register chan *Msg
	execute  chan *Msg
}
type Msg struct {
	From    string
	To      string
	Time    int64
	Content string
	Repeat  int
	Loop    int
}

// TODO Message
// 要 new 出一个环形的msg
func NewMessage() *Ring {
	return &Ring{
		Now:      make([]*Msg, 0),
		Dead:     make([]*Msg, 0),
		register: make(chan *Msg, 3000),
		execute:  make(chan *Msg, 3000),
	}
}
func (r *Ring) loop() {
	for {
		select {
		case m := <-r.register:
			fmt.Println(m)
			r.Now = append(r.Now, m)
			r.All++
			// r.execute <- m
		case m := <-r.execute:
			fmt.Println(m)
		default:

		}
	}
}

func Run() {
	r := NewMessage()
	go r.loop()
	rpc.Register(r) // 注册rpc服务
	lis, err := net.Listen("tcp", ":89")
	if err != nil {
		log.Fatalln("fatal error: ", err)
	}
	for {
		conn, err := lis.Accept() // 接收客户端连接请求
		if err != nil {
			continue
		}
		go func(conn net.Conn) { // 并发处理客户端请求
			fmt.Fprintf(os.Stdout, "%s", "new client in coming\n")
			jsonrpc.ServeConn(conn)
		}(conn)
	}
}
func (r *Ring) Push(msg *Msg, res *string) error {
	r.register <- msg
	*res = "ok"
	return nil
}
func (r *Ring) Pull(msg *Msg, res *string) error {
	return nil

}
