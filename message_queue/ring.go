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
	Now      *Msg
	Dead     []*Msg `用来存放没有成功推送出去的消息,消息都会立即推送出去，没推出去下方在Dead`
	All      int    `当前所有的消息数`
	mux      sync.RWMutex
	register chan *Msg
	head     chan *Msg
	execute  chan *Msg
}
type Msg struct {
	From    string
	To      string
	Time    int64
	Content string
	Repeat  int
	Loop    int
	Next    *Msg
}

// TODO Message
// 要 new 出一个环形的msg
func NewMessage() *Ring {
	return &Ring{
		Now:      nil,
		Dead:     make([]*Msg, 0),
		register: make(chan *Msg, 3000),
		execute:  make(chan *Msg, 3000),
		head:     make(chan *Msg, 3000),
	}
}
func (r *Ring) loop() {
	go func() { // 检测到入队后 直接把队列头部的取出来
		for {
			select {
			case <-r.head:
				//r.mux.RLock()
				//r.execute <- r.Now
				fmt.Println("head", r.Now)
				//r.Now = r.Now.Next
				//r.mux.RUnlock()
			}
		}
	}()

	for { // 检测入队
		select {
		case m := <-r.register:
			r.mux.RLock()
			if r.Now == nil {
				r.Now = m
				continue
			}
			p1 := r.Now
			p2 := r.Now.Next
			for p2 != nil {
				p1 = p2
				p2 = p1.Next
			}
			fmt.Println(p2)
			p2 = m
			r.mux.RUnlock()
			r.All++
			r.head <- m //放入执行队列
		}
	}
}

var RING *Ring

func Run() {
	RING := NewMessage()
	go RING.loop()
	_ = rpc.Register(RING) // 注册rpc服务
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
	x, ok := <-r.execute
	if ok {
		*res = x.Content
	} else {
		r.mux.RLock()
		r.Dead = append(r.Dead, x)
		r.mux.RUnlock()
		*res = "get error"
	}
	return nil

}
