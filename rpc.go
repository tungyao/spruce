package spruce

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
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

func (o *Operation) Get(args *OperationArgs, result *interface{}) error {
	log.Println("rpc get =>", args)
	*result = balala.Get(args.Key)
	return nil
}
func (o *Operation) Set(args *OperationArgs, result *int) error {
	log.Println("rpc set =>", args)
	*result = balala.Set(args.Key, args.Value, args.Expiration)
	return nil
}
func RpcStart(address string) {
	err := rpc.Register(new(Operation))
	if err != nil {
		log.Panicln(err)
	}
	rpc.HandleHTTP()
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("\n\nserver is listening =>", listen.Addr().String())
	err = http.Serve(listen, nil)
	if err != nil {
		log.Panicln(err)
	}
}
