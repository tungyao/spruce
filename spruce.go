package spruce

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-yaml/yaml"
	"google.golang.org/grpc"
)

const (
	FILE = iota
	MEMORY
)

type Config struct {
	ConfigType      int     `配置方式`
	DCSConfigFile   string  `分布式的配置文件路径`
	DNode           []DNode `采用内存的方式部署`
	Addr            string  `跑在哪个端口上`
	NowIP           string  `当前服务器运行的IP地址 暂时必须`
	KeepAlive       bool
	IsBackup        bool `自动备份`
	ConnChanBufSize int  `连接信道缓冲区大小`
	ConnChanMaxSize int  `最大连接数`
	MaxSlot         int  `最大hash槽数量`
}
type FileConfig struct {
	Config []DNode `yaml:"config"`
}

// type DCSConfig struct {
//	Name     string `json:"name"`
//	Ip       string `json:"ip"`
//	Weigh    int    `json:"weigh"`
//	Password string `json:"password"`
// }
type DNode struct {
	Name     string `yaml:"name"`
	Ip       string `yaml:"ip"`
	Weigh    int    `yaml:"weigh"`
	Password string `yaml:"password"`
}
type Slot struct {
	Count int
	Face  DNode
	Other []DNode
	cry   []uint
	All   []DNode
	Mux   sync.Mutex
}

var (
	AllSlot     []DNode
	EntrySlot   *Slot
	balala      *Hash
	randomMutex = sync.Mutex{}
	mux         sync.Mutex
	action      chan int // 1 是增加 -1 是减少 // 10个缓冲
)

// TODO Node
func setAllDNode(c []DNode) *Slot {
	cryptTable := make([]uint, 512)
	var (
		seed   uint = 0x00100001
		index1      = 0
		i           = 0
	)
	for index1 = 0; index1 < len(cryptTable); index1++ {
		var (
			tp1 uint
			tp2 uint
		)
		seed = (seed*uint(len(cryptTable)) + 3) % 0x2AAAAB
		tp1 = (seed & 0xFFFF) << 0x10
		seed = (seed*uint(len(cryptTable)) + 3) % 0x2AAAAB
		tp2 = seed & 0xFFFF
		cryptTable[index1] = tp1 | tp2
		i += 1
	}
	AllSlot = c
	return &Slot{Count: len(c), Face: c[0], Other: c[1:], All: c, cry: cryptTable}
}

// func New(config Config) *Slot {
//	CheckConfig(&config, Config{
//		ConfigType:    FILE,
//		DCSConfigFile: "./",
//		Addr:          "0.0.0.0:9102",
//		KeepAlive:     false,
//	})
//	return setAllDNode(ParseConfigFile(config.DCSConfigFile))
//
// }

// TODO Start
func StartSpruceDistributed(config Config) *Slot {
	balala = CreateHash(config.MaxSlot)
	// CheckConfig(&config, Config{
	//	ConfigType:      FILE,
	//	DCSConfigFile:   "./config.yml",
	//	Addr:            ":6998",
	//	ConnChanBufSize: 2048,
	//	ConnChanMaxSize: 2048,
	// })
	// region print logo
	fmt.Print(`
  ___ _ __  _ __ _   _  ___ ___ 
 / __| '_ \| '__| | | |/ __/ _ \
 \__ \ |_) | |  | |_| | (_|  __/
 |___/ .__/|_|   \__,_|\___\___|
     | |                        
     |_|                        
`)
	fmt.Print(`
Spruce is distributed key-value data based on go. 
Of course, we built it in an embedded way 
at the beginning of the design. You can also 
use ordinary map functions as easily `)
	// endregion
	switch config.ConfigType {
	case FILE:
		client(config)
	case MEMORY:
		createMemory(config)
	}
	return nil
}
func initDNode(p string) *Slot {
	d := ParseConfigFile(p)
	fmt.Print("\n\nrunning server\n")
	fmt.Print("id", "\t", "name", "\t", "ip", "\t", "weigh", "\n")
	for k, v := range d {
		fmt.Print(k, "\t", v.Name, "\t", v.Ip, "\t", v.Weigh, "\n")
	}
	return setAllDNode(d)
}

// TODO Client
func client(config Config) {
	EntrySlot = initDNode(config.DCSConfigFile)
	EntrySlot.Face.Ip = config.NowIP
	config.DNode = EntrySlot.All
	fmt.Println("now ip is ", config.NowIP, "we would contrast it")
	fmt.Println("server is running   =>", os.Getpid())
	// 启动RPC 要判断如果只有一台主机则不能启动RPC
	go NoRpcServer(&config)
	if EntrySlot.Count > 1 {
		RpcStart(config)
	}
	// 监听所有的slot
	// go listenAllSlotAction()
	// //  这里用来区分稳定的长连接
	// if config.KeepAlive {
	//	for {
	//		c, err := a.Accept()
	//		if err != nil {
	//			log.Println(err)
	//		}
	//		go func(co net.Conn) {
	//			tcpP, ok := co.(*net.TCPConn)
	//			if !ok {
	//				log.Println(ok)
	//				return
	//			}
	//			for {
	//				data := make([]byte, 1024)
	//				n, err := c.Read(data)
	//				//str := SplitString(data[:n], []byte("*$"))
	//				msg := make([]byte, 0)
	//				switch data[0] {
	//				case 0:
	//					msg = slot.Delete(data[:n])
	//				case 1:
	//					msg = slot.Set(data[:n])
	//					//return SendStatusMessage()
	//				case 2:
	//					msg = slot.Get(data[:n]).([]byte)
	//				case 3:
	//				case 4: // close this connection
	//					break
	//				}
	//				_, err = c.Write(msg)
	//				if err != nil {
	//					log.Println(err)
	//				}
	//			}
	//			tcpP.Close()
	//		}(c)
	//	}
	// }
	// // 应该写嵌入式了
	// for {
	//	c, err := a.Accept()
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	connChan <- c
	// }
	// }

}
func NoRpcServer(config *Config) {
	if config.KeepAlive {
		a, err := net.Listen("tcp", config.Addr)
		if err != nil {
			log.Println(549, err)
			return
		}
		defer a.Close()
		fmt.Println("\n\nserver is listening =>", a.Addr().String())
		fmt.Println("server is running   =>", os.Getpid())
		for {
			c, err := a.Accept()
			if err != nil {
				log.Println(216, err)
			}
			go memoryServeHandle(c)
		}
	}
	a, err := net.Listen("tcp", config.Addr)
	if err != nil {
		log.Println(205, err)
		return
	}
	defer a.Close()
	connChan := make(chan net.Conn, config.ConnChanBufSize)
	connChanSize := make(chan int)
	connSize := 0
	go func() {
		for cc := range connChanSize {
			connSize += cc
		}
	}()
	for i := 0; i < config.ConnChanMaxSize; i++ {
		go func() {
			for conn := range connChan {
				connChanSize <- 1
				EchoNoKeepAlive(conn, EntrySlot)
				connChanSize <- 1
			}
		}()
	}
}
func EchoNoKeepAlive(c net.Conn, slot *Slot) {
	defer func() {
		x := recover()
		log.Println("223 line", x)
		c.Close()
	}()
	data := make([]byte, 1024)
	n, err := c.Read(data)
	// str := SplitString(data[:n], []byte("*$"))
	msg := make([]byte, 0)
	switch data[0] {
	case 0:
		msg = slot.Delete(data[:n])
	case 1:
		msg = slot.Set(data[:n])
		// return SendStatusMessage()
	case 2:
		getValue := slot.Get(data[:n])
		if getValue == nil {
			msg = nil
		} else {
			msg = getValue.([]byte)
		}
	case 3:
		// storage  into the current path
		slot.Storage()
	}
	_, err = c.Write(msg)
	if err != nil {
		log.Println(256, err)
	}
}
func (s *Slot) Storage() {
	balala.Storage()
}

func (s *Slot) Get(lang []byte) interface{} {
	n := s.getHashPos(lang[11:])
	fmt.Println("get value of", s.All[n].Ip, "slot", string(lang[11:]))
	if s.All[n].Ip == s.Face.Ip {
		return balala.Get(lang[11:])
	} else {
		return GetRpc(&OperationArgs{Key: lang[11:]}, s.All[n].Ip)
	}
}
func (s *Slot) Set(lang []byte) []byte {
	key, value := SplitKeyValue(lang[11:])
	ns := s.getHashPos(key)
	fmt.Println("set value to", s.All[ns].Ip, "slot", string(key))
	ti := ParsingExpirationDate(lang[2:4]).(int64)
	if s.All[ns].Ip == s.Face.Ip {
		it := balala.Set(key, value, ti)
		fmt.Println("saved", it)
		return []byte(strconv.Itoa(it))
	} else {
		return []byte{uint8(SetRpc(&OperationArgs{
			Key:        key,
			Value:      value,
			Expiration: ti,
		}, s.All[ns].Ip))}
	}
}
func GetRpc(args *OperationArgs, address string) []byte {
	conn, err := grpc.Dial(address)
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()
	op := NewOperationClient(conn)
	ctx, cancle := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancle()
	res, err := op.Get(ctx, args)
	if err != nil {
		log.Panicln(err)
	}
	return res.Value
}
func DeleteRpc(args *OperationArgs, address string) []byte {
	conn, err := grpc.Dial(address)
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()
	op := NewOperationClient(conn)
	ctx, cancle := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancle()
	res, err := op.Delete(ctx, args)
	if err != nil {
		log.Panicln(err)
	}
	return res.Value
}
func SetRpc(args *OperationArgs, address string) int {
	client, err := grpc.Dial(address)
	if err != nil {
		return 0
	}
	defer client.Close()
	nrc := NewOperationClient(client)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rest, err := nrc.Set(ctx, args)
	if err != nil {
		log.Panicln("343", err)
	}
	return int(rest.Position)
}
func getRemote(lang []byte, ip string) []byte {
	con, err := net.Dial("tcp", ip)
	if err != nil {
		log.Println(320, err)
		return nil
	}
	defer con.Close()
	con.Write(lang)
	return GetData(con)
}
func (s *Slot) Position(key []byte) int {
	return int(s.getHashPos(key))
}
func (s *Slot) Delete(lang []byte) []byte {
	key := lang[11:]
	if len(key) < 1 {
		return nil
	}
	ns := s.getHashPos(key)
	fmt.Println("delete value of", s.Face.Ip, "slot", string(key))
	if s.All[ns].Ip == s.Face.Ip {
		return balala.Delete(key)
	} else {
		return DeleteRpc(&OperationArgs{Key: lang[11:]}, s.All[ns].Ip)
	}
}
func GetData(a net.Conn) []byte {
	out := make([][]byte, 0)
	o := make([]byte, 0)
	for {
		data := make([]byte, 1024)
		n, err := a.Read(data)
		out = append(out, data[:n])
		if n == 0 || err == io.EOF {
			break
		}
	}
	for _, v := range out {
		for _, j := range v {
			if j == 0 {
				continue
			}
			o = append(o, j)
		}
	}
	return o
}
func (s *Slot) hashString(str []byte, hashcode uint) uint {
	var (
		key        = str
		seed1 uint = 0x7FED7FED
		seed2 uint = 0xEEEEEEEE
		ch    uint = 0
	)
	for _, v := range key {
		ch = uint(v)
		seed1 = s.cry[(hashcode<<2)+ch] ^ (seed1 + seed2)
		seed2 = ch + seed1 + seed2 + (seed2 << 5) + 3
	}
	return seed1
}
func (s *Slot) getHashPos(str []byte) uint {
	var (
		hOffset  uint = 0
		hashA    uint = 1
		hashB    uint = 2
		nHash    uint = s.hashString(str, hOffset)
		nHashA   uint = s.hashString(str, hashA)
		nHashB   uint = s.hashString(str, hashB)
		nHashPos uint = (nHash + nHashA + nHashB) % uint(s.Count)
	)
	return nHashPos
}
func ParseConfigFile(path string) []DNode {
	// fmt.Println(path)
	// yml := ymload.Format(path)
	var config FileConfig
	fs, err := os.Open(path)
	if err != nil {
		log.Panicln(err)
	}
	b, _ := ioutil.ReadAll(fs)
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Panicln(err)
	}
	// dn := make([]DNode, 0)
	// for _, v := range config.Config {
	//	ds := DNode{}
	//	ds.Ip = v.Ip
	//	ds.Name = v.Name
	//	ds.Weigh = v.Weigh
	//	ds.Pwd = v.Password
	//	dn = append(dn, ds)
	// }
	return config.Config
	// dn := make([]DNode, 0)
	// for _, v := range yml {
	//	ds := DNode{}
	//	if v["ip"] != nil {
	//		ds.Ip = v["ip"].(string)
	//	}
	//	if v["name"] != nil {
	//		ds.Name = v["name"].(string)
	//	}
	//	if v["weight"] != nil {
	//		i, _ := strconv.Atoi(v["weight"].(string))
	//		ds.Weigh = i
	//	}
	//	if v["password"] != nil {
	//		ds.Pwd = v["password"].(string)
	//	}
	//	dn = append(dn, ds)
	// }
	// return dn
	// f, err := os.Open(path)
	// defer f.Close()
	// if err != nil {
	//	log.Println("open config file error", err)
	// }
	// stat, _ := f.Stat()
	// get := make([]byte, stat.Size())
	// _, err = f.Read(get)
	// if err != nil {
	//	log.Panic(err)
	// }
	// isgroup := false
	// str := make([]byte, 0)
	// dn := make([]DNode, 0)
	// for i := 0; i < len(get); i++ {
	//	if get[i] == 32 {
	//		continue
	//	}
	//	if get[i] == 123 {
	//		isgroup = true
	//	}
	//	if get[i] == 125 {
	//		isgroup = false
	//	}
	//	if isgroup && get[i] != 123 {
	//		str = append(str, get[i])
	//	}
	// }
	// group := SplitString(str, []byte("\n\n"))
	// // 到这一部可以开始解析数据到出来
	// for _, v := range group {
	//	ds := DNode{}
	//	column := SplitString(v, []byte("\n"))
	//	for _, j := range column {
	//		name := FindString(j, []byte("name="))
	//		if name != nil {
	//			ds.Name = string(name.([]uint8))
	//		}
	//		ip := FindString(j, []byte("ip="))
	//		if ip != nil {
	//			ds.Ip = string(ip.([]uint8))
	//		}
	//		password := FindString(j, []byte("password="))
	//		if password != nil {
	//			ds.Pwd = string(password.([]uint8))
	//		}
	//		weight := FindString(j, []byte("weight="))
	//		if weight != nil {
	//			s := weight.([]uint8)
	//			str := make([]byte, 0)
	//			for _, v := range s {
	//				if v <= 57 && v >= 48 {
	//					str = append(str, v)
	//				}
	//			}
	//			d, err := strconv.Atoi(string(str))
	//			if err != nil {
	//				log.Panicln(err)
	//			}
	//			ds.Weigh = d
	//		}
	//	}
	//	dn = append(dn, ds)
	// }
	// return dn
}

// 增加新的插槽
func AddSlot() {

}

// 删除新的插槽
func DropSlot() {

}

// 这个方法是用于重置slot ，我们要重新计算 hash槽
// 实现方法 ,将每台电脑的数组取出来，重新取余，将值转移到对应slot，本机删除，如果计算结果是本机，那么不转移

func (s *Slot) ResetSlot(n DNode) {
	// 首先先把slot给锁住，防止出现混乱
	s.Mux.Lock()
	alls := s.All
	for _, v := range alls {
		getRemote([]byte("reset*$"+n.Name+"*$"+n.Ip+"*$"+strconv.Itoa(n.Weigh)+"*$"+n.Password), v.Ip)
	}
	s.Mux.Unlock()
}

// 检测slot的变化
func listenAllSlotAction() {
	log.Println("start listening SlotAction")
	for at := range action {
		if at == 1 {

		}
	}
}

// memory control 通过嵌入式来配置spruce
func createMemory(config Config) {
	fmt.Print("\n\nrunning server\n")
	fmt.Print("id", "\t", "name", "\t", "ip", "\t", "weigh", "\n")
	for k, v := range config.DNode {
		fmt.Print(k, "\t", v.Name, "\t", v.Ip, "\t", v.Weigh, "\n")
	}
	EntrySlot = setAllDNode(config.DNode)
	EntrySlot.Face.Ip = config.NowIP
	if len(config.DNode) > 1 {
		go RpcStart(config)
	}
	// if slot.Count <= 1 {
	createMemoryServe(config, EntrySlot)
	// } else { //大于一台则只能只用RPC

	// }

}
func createMemoryServe(config Config, s *Slot) {
	if config.KeepAlive {
		a, err := net.Listen("tcp", config.Addr)
		if err != nil {
			log.Println(549, err)
			return
		}
		defer a.Close()
		fmt.Println("\n\nserver is listening =>", a.Addr().String())
		fmt.Println("server is running   =>", os.Getpid())
		for {
			c, err := a.Accept()
			if err != nil {
				log.Println(576, err)
			}
			go memoryServeHandle(c)
		}
	}
	a, err := net.Listen("tcp", config.Addr)
	if err != nil {
		log.Println(566, err)
		return
	}
	fmt.Println("now ip is ", config.NowIP, "we would contrast it")
	fmt.Println("\n\nserver is listening =>", a.Addr().String())
	fmt.Println("server is running   =>", os.Getpid())
	go listenAllSlotAction()
	// 应该写嵌入式了
	for {
		c, err := a.Accept()
		if err != nil {
			log.Println("601", err)
		}
		go memoryServeHandleConn(c)
	}
}
func memoryServeHandleConn(c net.Conn) {
	defer c.Close()
	data := make([]byte, 2048)
	n, err := c.Read(data)
	// log.Println("get bytes", data[:n], n)
	if n < 11 {
		log.Println(588, err)
		return
	}
	msg := make([]byte, 0)
	switch data[0] {
	case 0:
		msg = EntrySlot.Delete(data[:n])
	case 1:
		msg = EntrySlot.Set(data[:n])
		// return SendStatusMessage()
	case 2:
		if m := EntrySlot.Get(data[:n]); m == nil {
			msg = []byte{}
		} else {
			msg = m.([]byte)
		}
	case 3:
		msg = CreateUUID(int(data[1]), data[11:n], CreateNewId(4))
	}
	_, err = c.Write(msg)
	if err != nil {
		log.Println("628", err)
	}
}
func memoryServeHandle(c net.Conn) {
	for {
		data := make([]byte, 2048)
		n, err := c.Read(data)
		if err != nil {
			fmt.Println(635, err)
			err = c.Close()
			break
		}
		msg := []byte{0}
		if n <= 11 {
			goto end
		}
		switch data[0] {
		case 0:
			msg = EntrySlot.Delete(data[:n])
		case 1:
			msg = EntrySlot.Set(data[:n])
		case 2:
			if m := EntrySlot.Get(data[:n]); m == nil {
			} else {
				msg = m.([]byte)
			}
		case 3:
		}
		goto end
	end:
		_, err = c.Write(msg)
		if err != nil {
			err = c.Close()
			log.Println("656", err)
			break
		}
	}
}
