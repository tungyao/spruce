package spruce

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/rpc"
	"time"
)

// rpc method

type Operation struct {
	UnimplementedOperationServer
}

func (o *Operation) Get(ctx context.Context, in *OperationArgs) (*Result, error) {
	log.Println("rpc get =>", in.String())
	result := &Result{Value: balala.Get(in.Key).([]byte)}
	return result, nil
}
func (o *Operation) Delete(ctx context.Context, in *OperationArgs) (*DeleteResult, error) {
	log.Println("rpc set =>", in.String())
	result := &DeleteResult{Value: balala.Delete(in.Key)}
	return result, nil
}
func (o *Operation) Set(ctx context.Context, in *OperationArgs) (*SetResult, error) {
	log.Println("rpc set =>", in.String())
	result := &SetResult{Position: int64(balala.Set(in.Key, in.Value, in.Expiration))}
	return result, nil
}

// 长连接心跳检测

type Watcher struct {
	UnimplementedWatcherServer
}

//func (w *Watcher) Ping(ip string) int8 {
//
//}
func (w *Watcher) Pong(ctx context.Context, in *WatcherData) (*WatcherResult, error) {
	result := &WatcherResult{Res: 12}
	return result, nil
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
	lis, err := net.Listen("tcp", ":6999")
	if err != nil {
		log.Panicln(err)
	}
	s := grpc.NewServer()
	RegisterOperationServer(s, &Operation{})
	RegisterWatcherServer(s, &Watcher{})
	log.Println("\n\ngRPC is running")
	// go startWatcher(config)
	s.Serve(lis)
}
