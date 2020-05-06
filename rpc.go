package spruce

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

// rpc method
//
type Operation struct {
}
type OperationArgs struct {
	Key        []byte
	Value      interface{}
	Expiration int64
}

// 心跳检测

func (o *Operation) Get(args *OperationArgs, result *interface{}) error {
	log.Println("rpc get =>", args)
	*result = balala.Get(args.Key)
	return nil
}
func (o *Operation) Delete(args *OperationArgs, result *interface{}) error {
	log.Println("rpc get =>", args)
	*result = balala.Delete(args.Key)
	return nil
}
func (o *Operation) Set(args *OperationArgs, result *int) error {
	log.Println("rpc set =>", args)
	*result = balala.Set(args.Key, args.Value, args.Expiration)
	return nil
}

type Watcher struct {
}
type WatcherData struct {
	Time int64
}

//func (w *Watcher) Ping(ip string) int8 {
//
//}
func (w *Watcher) Pong(args *WatcherData, result *int8) error {
	var x int8 = 12
	*result = x
	return nil
}
func (w *Watcher) Do(args *WatcherData, result *int8) error {
	var x int8 = 13
	*result = x
	return nil
}
func (w *Watcher) Dead(args *WatcherData, result *int8) error {
	var x int8 = 14
	*result = x
	return nil
}
func startWatcher(dsc Config) {
	log.Println("starting rpc watcher ...")
	log.Print(`
__     __     ______     ______   ______     __  __     ______     ______
/\ \  _ \ \   /\  __ \   /\__  _\ /\  ___\   /\ \_\ \   /\  ___\   /\  == \
\ \ \/ ".\ \  \ \  __ \  \/_/\ \/ \ \ \____  \ \  __ \  \ \  __\   \ \  __<
\ \__/".~\_\  \ \_\ \_\    \ \_\  \ \_____\  \ \_\ \_\  \ \_____\  \ \_\ \_\
 \/_/   \/_/   \/_/\/_/     \/_/   \/_____/   \/_/\/_/   \/_____/   \/_/ /_/

`)
	for _, v := range dsc.DNode {
		if v.Ip == dsc.NowIP {
			continue
		}
		c, err := rpc.Dial("tcp", v.Ip)
		for err != nil {
			c, err = rpc.Dial("tcp", v.Ip)
			<-time.After(time.Second * 2)
		}
		var x int8
		err = c.Call("Watcher.Pong", &WatcherData{}, &x)
		if err != nil {
			log.Println(err)
		}
		log.Println("get ping data =>", x)
		c.Close()
	}
}
func RpcStart(config Config) {
	err := rpc.Register(new(Operation))
	err = rpc.Register(new(Watcher))
	if err != nil {
		log.Panicln(err)
	}
	listen, err := net.Listen("tcp", ":6999")
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("\n\nRPC is listening =>", listen.Addr().String())
	go startWatcher(config)
	//rpc.NewServer()
	rpc.Accept(listen)
}
