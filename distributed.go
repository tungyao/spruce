package spruce

import (
	"fmt"
	"github.com/tungyao/ymload"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	FILE = iota
	MEMORY
)

type Config struct {
	ConfigType    int         `配置方式`
	DCSConfigFile string      `分布式的配置文件路径`
	DCSConfigs    []DCSConfig `采用内存的方式部署`
	Addr          string      `跑在哪个端口上`
	NowIP         string      `当前服务器运行的IP地址 暂时必须`
	KeepAlive     bool
	IsBackup      bool
}
type DCSConfig struct {
	Name     string
	Ip       string
	Weigh    int
	Password string
}
type DNode struct {
	Name  string
	IP    string
	Weigh int
	Pwd   string
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
	AllSlot      []DNode
	slot         int
	balala       = CreateHash(4096)
	randomMutex  = sync.Mutex{}
	mux          sync.Mutex
	action       chan int // 1 是增加 -1 是减少 // 10个缓冲
	nMessageRing *MessageRing
	pMsg         chan msg
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
func New(config Config) *Slot {
	CheckConfig(&config, Config{
		ConfigType:    FILE,
		DCSConfigFile: "./",
		Addr:          "0.0.0.0:9102",
		KeepAlive:     false,
	})
	return setAllDNode(ParseConfigFile(config.DCSConfigFile))

}

// TODO Start
func StartSpruceDistributed(config Config) *Slot {
	CheckConfig(&config, Config{
		ConfigType:    FILE,
		DCSConfigFile: "./spruce.yml",
		Addr:          ":9102",
	})
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
		return createMemory(config)
	}
	return nil
}
func initDNode(p string) *Slot {
	d := ParseConfigFile(p)
	fmt.Print("\n\nrunning server\n")
	fmt.Print("id", "\t", "name", "\t", "ip", "\t", "weigh", "\n")
	for k, v := range d {
		fmt.Print(k, "\t", v.Name, "\t", v.IP, "\t", v.Weigh, "\n")
	}
	return setAllDNode(d)
}

// TODO Client
func client(config Config) {
	slot := initDNode(config.DCSConfigFile)
	slot.Face.IP = config.NowIP

	a, err := net.Listen("tcp", config.Addr)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("now ip is ", config.NowIP, "we would contrast it")
	fmt.Println("\n\nserver is listening =>", a.Addr().String())
	fmt.Println("server is running   =>", os.Getpid())
	go listenAllSlotAction()
	//  这里用来区分稳定的长连接
	if config.KeepAlive {
		for {
			c, err := a.Accept()
			if err != nil {
				log.Println(err)
			}
			go func(co net.Conn) {
				tcpP, ok := co.(*net.TCPConn)
				if !ok {
					log.Println(ok)
					return
				}
				for {
					data := make([]byte, 1024)
					n, err := c.Read(data)
					//str := SplitString(data[:n], []byte("*$"))
					msg := make([]byte, 0)
					switch data[0] {
					case 0:
						msg = slot.Delete(data[:n])
					case 1:
						msg = slot.Set(data[:n])
						//return SendStatusMessage()
					case 2:
						msg = slot.Get(data[:n])
					case 3:
					case 4: // close this connection
						break
					}
					_, err = c.Write(msg)
					if err != nil {
						log.Println(err)
					}
				}
				tcpP.Close()
			}(c)
		}
	}
	// 应该写嵌入式了
	for {
		c, err := a.Accept()
		if err != nil {
			log.Println(err)
		}
		go func(c net.Conn) {
			data := make([]byte, 1024)
			n, err := c.Read(data)
			//str := SplitString(data[:n], []byte("*$"))
			msg := make([]byte, 0)
			switch data[0] {
			case 0:
				msg = slot.Delete(data[:n])
			case 1:
				msg = slot.Set(data[:n])
				//return SendStatusMessage()
			case 2:
				msg = slot.Get(data[:n])
			case 3:
			}
			_, err = c.Write(msg)
			if err != nil {
				log.Println(err)
			}
			err = c.Close()
		}(c)

	}
	//}

}
func (s *Slot) Get(lang []byte) []byte {
	n := s.getHashPos(lang[11:])
	fmt.Println("get value of", n, "slot", string(lang[11:]))
	if s.All[n].IP == s.Face.IP {
		return balala.Get(lang[11:])
	} else {
		return getRemote(SendGetMessage(lang[11:]), s.All[n].IP)
	}
}
func (s *Slot) Position(key []byte) int {
	return int(s.getHashPos(key))
}
func (s *Slot) Set(lang []byte) []byte {
	key, value := SplitKeyValue(lang[11:])
	ns := s.getHashPos(key)
	fmt.Println("set value to", s.Face.IP, "slot", string(key))
	if s.All[ns].IP == s.Face.IP {
		fmt.Println("save")
		ti := ParsingExpirationDate(lang[2:4]).(int64)
		it := balala.Set(key, value, ti)
		return []byte{uint8(it)}
	} else {
		return getRemote(lang, s.All[ns].IP)
	}
}
func (s *Slot) Delete(lang []byte) []byte {
	key := lang[11:]
	if len(key) < 1 {
		return nil
	}
	ns := s.getHashPos(key)
	fmt.Println("delete value of", s.Face.IP, "slot", string(key))
	if s.All[ns].IP == s.Face.IP {
		return balala.Delete(key)
	} else {
		return getRemote(lang, s.All[ns].IP)
	}
}
func getRemote(lang []byte, ip string) []byte {
	con, err := net.Dial("tcp", ip)
	if err != nil {
		log.Println(295, err)
		return nil
	}
	defer con.Close()
	con.Write(lang)
	return GetData(con)
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
	fmt.Println(path)
	yml := ymload.Format(path)
	dn := make([]DNode, 0)
	for _, v := range yml {
		ds := DNode{}
		if v["ip"] != nil {
			ds.IP = v["ip"].(string)
		}
		if v["name"] != nil {
			ds.Name = v["name"].(string)
		}
		if v["weight"] != nil {
			i, _ := strconv.Atoi(v["weight"].(string))
			ds.Weigh = i
		}
		if v["password"] != nil {
			ds.Pwd = v["password"].(string)
		}
		dn = append(dn, ds)
	}
	return dn
	//f, err := os.Open(path)
	//defer f.Close()
	//if err != nil {
	//	log.Println("open config file error", err)
	//}
	//stat, _ := f.Stat()
	//get := make([]byte, stat.Size())
	//_, err = f.Read(get)
	//if err != nil {
	//	log.Panic(err)
	//}
	//isgroup := false
	//str := make([]byte, 0)
	//dn := make([]DNode, 0)
	//for i := 0; i < len(get); i++ {
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
	//}
	//group := SplitString(str, []byte("\n\n"))
	//// 到这一部可以开始解析数据到出来
	//for _, v := range group {
	//	ds := DNode{}
	//	column := SplitString(v, []byte("\n"))
	//	for _, j := range column {
	//		name := FindString(j, []byte("name="))
	//		if name != nil {
	//			ds.Name = string(name.([]uint8))
	//		}
	//		ip := FindString(j, []byte("ip="))
	//		if ip != nil {
	//			ds.IP = string(ip.([]uint8))
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
	//}
	//return dn
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
		getRemote([]byte("reset*$"+n.Name+"*$"+n.IP+"*$"+strconv.Itoa(n.Weigh)+"*$"+n.Pwd), v.IP)
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
func getRandomInt(start, end int) int {
	randomMutex.Lock()
	<-time.After(1 * time.Nanosecond)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := start + r.Intn(end-start+1)
	randomMutex.Unlock()
	return n
}

// memory control 通过嵌入式来配置spruce
func createMemory(config Config) *Slot {
	d := make([]DNode, len(config.DCSConfigs))
	fmt.Print("\n\nrunning server\n")
	fmt.Print("id", "\t", "name", "\t", "ip", "\t", "weigh", "\n")
	for k, v := range config.DCSConfigs {
		d[k] = DNode{
			Name:  v.Name,
			IP:    v.Ip,
			Weigh: v.Weigh,
			Pwd:   v.Password,
		}
		fmt.Print(k, "\t", v.Name, "\t", v.Ip, "\t", v.Weigh, "\n")
	}
	slot := setAllDNode(d)
	slot.Face.IP = config.NowIP
	go createMemoryServe(config, slot)
	return slot
}
func createMemoryServe(config Config, s *Slot) {
	if config.KeepAlive {
		tcpAddr, err := net.ResolveTCPAddr("tcp", config.Addr) //创建 tcpAddr数据
		a, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			log.Println(err)
			return
		}
		defer a.Close()
		fmt.Println("\n\nserver is listening =>", a.Addr().String())
		fmt.Println("server is running   =>", os.Getpid())
		for {
			c, err := a.AcceptTCP()
			fmt.Println(123)
			if err != nil {
				log.Println(err)
			}
			go memoryServeHandle(c, s)
		}
	}
	a, err := net.Listen("tcp", config.Addr)
	if err != nil {
		log.Println(err)
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
			log.Println(err)
		}
		go memoryServeHandleConn(c, s)
	}
}
func memoryServeHandleConn(c net.Conn, slot *Slot) {
	data := make([]byte, 1024)
	n, err := c.Read(data)
	log.Println("get bytes", data[:n], n)
	if err != nil {
		log.Println(err)
	} else {
		err = c.Close()
	}
	msg := make([]byte, 0)
	switch data[0] {
	case 0:
		msg = slot.Delete(data[:n])
	case 1:
		msg = slot.Set(data[:n])
		//return SendStatusMessage()
	case 2:
		msg = slot.Get(data[:n])
	case 3:
	}
	_, err = c.Write(msg)
	if err != nil {
		log.Println(err)
	}
	err = c.Close()
}
func memoryServeHandle(c *net.TCPConn, slot *Slot) {
	data := make([]byte, 1024)
	n, err := c.Read(data)
	log.Println("get bytes", string(data[:n]), n)
	if err != nil {
		log.Println(err)
	} else {
		err = c.CloseRead()
	}
	msg := make([]byte, 0)
	if n <= 11 {
		goto end
	}
	switch data[0] {
	case 0:
		msg = slot.Delete(data[:n])
	case 1:
		msg = slot.Set(data[:n])
		//return SendStatusMessage()
	case 2:
		msg = slot.Get(data[:n])
	case 3:
	}
end:
	_, err = c.Write(msg)
	if err != nil {
		log.Println(err)
	}
	err = c.Close()
}
