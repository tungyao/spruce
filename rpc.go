package spruce

import (
	"fmt"
	awesome "git.yaop.ink/tungyao/awesome-pool"
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
	T []*awesome.Pool
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
func startWatcher(dsc *[]DNode) {
	log.Println("starting rpc watcher ...")
	log.Print(`
__     __     ______     ______   ______     __  __     ______     ______
/\ \  _ \ \   /\  __ \   /\__  _\ /\  ___\   /\ \_\ \   /\  ___\   /\  == \
\ \ \/ ".\ \  \ \  __ \  \/_/\ \/ \ \ \____  \ \  __ \  \ \  __\   \ \  __<
\ \__/".~\_\  \ \_\ \_\    \ \_\  \ \_____\  \ \_\ \_\  \ \_____\  \ \_\ \_\
 \/_/   \/_/   \/_/\/_/     \/_/   \/_____/   \/_/\/_/   \/_____/   \/_/ /_/

`)
	wc := new(Watcher)
	wc.T = make([]*awesome.Pool, len(*dsc))
	goto Restart
Restart:
	var errx error
	for k, v := range *dsc {
		log.Println("ping address =>", v.Ip)
		wc.T[k], errx = awesome.NewPool(10, v.Ip)
		if errx != nil {
			log.Println("ready reconnection ......")
			<-time.After(time.Second * 5)
			goto Restart
		}
	}

	for {
		log.Println("monitor the watcher")
		for _, v := range wc.T {
			c := v.Get()
			if c == nil {
				log.Println("ready reconnection ......")
				<-time.After(time.Second * 5)
				goto Restart
			}
			client := rpc.NewClient(c.Conn)
			var x int8
			err := client.Call("Watcher.Pong", &WatcherData{}, &x)
			if err != nil {
				log.Println(err)
			}
			log.Println("get ping data =>", x)
			client.Close()
		}
		<-time.After(time.Second * 2)
	}
}
func RpcStart(address Config) {
	newDCS := make([]DNode, 0)
	for _, v := range address.DNode {
		if v.Ip != address.NowIP {
			newDCS = append(newDCS, v)
		}
	}
	err := rpc.Register(new(Operation))
	err = rpc.Register(new(Watcher))
	if err != nil {
		log.Panicln(err)
	}
	listen, err := net.Listen("tcp", ":82")
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("\n\nRPC is listening =>", listen.Addr().String())
	go startWatcher(&newDCS)
	//rpc.NewServer()
	rpc.Accept(listen)
}
