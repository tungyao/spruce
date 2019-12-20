package main

import (
	"./spruce"
	"flag"
	"fmt"
	"os"
)

var (
	use    string
	path   string
	create bool
	start  bool
	stop   bool
	addr   string
	keep   bool
	nowip string
)

func init() {
	flag.StringVar(&use, "use", "file", "help message for flag name")
	flag.StringVar(&nowip, "nowip", "0.0.0.0", "we would contrast it")
	flag.StringVar(&path, "path", "./spruce.yml", "config file path ,default current folder")
	flag.BoolVar(&create, "create", false, "create new config file")
	flag.BoolVar(&start, "start", false, "start gate-way")
	flag.BoolVar(&stop, "stop", false, "stop gate-way")
	flag.BoolVar(&keep, "keep", false, "keep alive connect")
	flag.StringVar(&addr, "addr", "0.0.0.0:9102", "listen port")
}
func main() {
	flag.Parse()
	if create {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 777)
		if err != nil {
			fmt.Println(30, err)
			f.Close()
			os.Exit(0)
		}
		_, err = f.Write([]byte("main_server:\r  name: client0\r  ip: 123123\r  password: 123\rtwo_server:\r  name: client1\r  ip: 123123\r  password: 123"))
		err = f.Close()
		if err != nil {
			fmt.Println(err)
		}
			os.Exit(0)
	}

	spruce.StartSpruceDistributed(spruce.Config{
		ConfigType: spruce.FILE,
		DCSConfigFile: "./spruce.yml",
		Addr:       addr,
		NowIP:      nowip,
		KeepAlive:  keep,
	})
}
