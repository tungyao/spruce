package spruce

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	AllSlot     []DNode
	slot        int
	balala      = CreateHash(4096)
	randomMutex = sync.Mutex{}
)

const (
	FILE = iota
	COMMAND
)

type Config struct {
	ConfigType    int    `配置方式`
	DCSConfigFile string `分布式的配置文件路径`
	Addr          string
	KeepAlive     bool
	IsBackup      bool
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
}

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
func StartSpruceDistributed(config Config) {
	CheckConfig(&config, Config{
		ConfigType:    FILE,
		DCSConfigFile: "./spruce.spe",
		Addr:          ":9102",
	})
	fmt.Println(config)
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
	client(config)
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
func client(config Config) {
	slot := initDNode(config.DCSConfigFile)
	//if config.KeepAlive {
	//	tcpAddr, err := net.ResolveTCPAddr("tcp", config.Addr) //创建 tcpAddr数据
	//	a, err := net.ListenTCP("tcp", tcpAddr)
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//	slot.Face.IP = config.Addr
	//	fmt.Println("\n\nserver is listening =>", a.Addr().String())
	//	fmt.Println("server is running   =>", os.Getpid())
	//	for {
	//		c, err := a.AcceptTCP()
	//		if err !=nil {
	//		    log.Println(err)
	//			_ = c.Close()
	//		}
	//		go func(c *net.TCPConn) {
	//			err = c.SetKeepAlive(true)
	//			if err != nil {
	//				log.Println(err)
	//				err  = c.Close()
	//				if err !=nil {
	//					log.Println(err)
	//				}
	//			}
	//			data := make([]byte, 1024)
	//			n, err := c.Read(data)
	//			if err !=nil {
	//				log.Println(err)
	//			}else{
	//				err = c.CloseRead()
	//			}
	//			fmt.Println(string(data))
	//			str := SplitString(data[:n], []byte("*$"))
	//			msg := ""
	//			switch string(str[0]) {
	//			case "get":
	//				msg = slot.Get(string(str[1]))
	//			case "set":
	//				if len(str) == 3 {
	//					slot.Set(string(str[1]), string(str[2]), 0)
	//					msg = "1"
	//				} else if len(str) == 4 {
	//					ns, err := strconv.Atoi(string(str[3]))
	//					if err == nil {
	//						slot.Set(string(str[1]), string(str[2]), ns)
	//						msg = "1"
	//					} else {
	//						msg = ""
	//					}
	//				}
	//			default:
	//				msg = ""
	//			}
	//			_, err = c.Write([]byte(msg))
	//			if err != nil {
	//				log.Println(err)
	//			}else{
	//				err = c.CloseWrite()
	//			}
	//			//err = c.Close()
	//			//if err !=nil {
	//			//    log.Println(err)
	//			//}
	//		}(c)
	//
	//	}
	//} else {
	//tcpAddr, err := net.ResolveTCPAddr("tcp", config.Addr) //创建 tcpAddr数据
	a, err := net.Listen("tcp", config.Addr)
	if err != nil {
		log.Println(err)
		return
	}
	slot.Face.IP = config.Addr
	fmt.Println("\n\nserver is listening =>", a.Addr().String())
	fmt.Println("server is running   =>", os.Getpid())
	for {
		c, err := a.Accept()

		if err != nil {
			log.Println(err)
		}
		go func(c net.Conn) {
			data := make([]byte, 1024)
			n, err := c.Read(data)
			str := SplitString(data[:n], []byte("*$"))
			msg := ""
			switch string(str[0]) {
			case "get":
				msg = slot.Get(string(str[1]))
			case "set":
				if len(str) == 3 {
					slot.Set(string(str[1]), string(str[2]), 0)
					msg = "1"
				} else if len(str) == 4 {
					ns, err := strconv.Atoi(string(str[3]))
					if err == nil {
						slot.Set(string(str[1]), string(str[2]), ns)
						msg = "1"
					} else {
						msg = ""
					}
				}
			default:
				msg = ""
			}
			_, err = c.Write([]byte(msg))
			if err != nil {
				log.Println(err)
			}
			err = c.Close()
		}(c)

	}
	//}

}
func (s *Slot) Get(key string) string {
	n := s.getHashPos([]rune(key))
	if s.All[n].IP == s.Face.IP {
		return balala.Get(key)
	} else {
		return getRemote([]byte("get*$"+key), s.All[n].IP)
	}
}
func (s *Slot) Set(n ...interface{}) string {

	key := n[0].(string)
	value := n[1].(string)
	ns := s.getHashPos([]rune(key))
	if s.All[ns].IP == s.Face.IP {
		return string(balala.Set(key, value, int64(n[2].(int))))
	} else {
		return getRemote([]byte("set*$"+key+"*$"+value+"*$"+strconv.Itoa(n[2].(int))), s.All[ns].IP)
	}
}
func getRemote(lang []byte, ip string) string {
	con, err := net.Dial("tcp", ip)
	if err != nil {
		log.Println(252, err)
		return ""
	}
	defer con.Close()
	con.Write(lang)
	return string(GetData(con))
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
func (s *Slot) hashString(str []rune, hashcode uint) uint {
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
func (s *Slot) getHashPos(str []rune) uint {
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
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Println("open config file error", err)
	}
	stat, _ := f.Stat()
	get := make([]byte, stat.Size())
	_, err = f.Read(get)
	if err != nil {
		log.Panic(err)
	}
	isgroup := false
	str := make([]byte, 0)
	dn := make([]DNode, 0)
	for i := 0; i < len(get); i++ {
		if get[i] == 32 {
			continue
		}
		if get[i] == 123 {
			isgroup = true
		}
		if get[i] == 125 {
			isgroup = false
		}
		if isgroup && get[i] != 123 {
			str = append(str, get[i])
		}
	}
	group := SplitString(str, []byte("\n\n"))
	// 到这一部可以开始解析数据到出来
	for _, v := range group {
		ds := DNode{}
		column := SplitString(v, []byte("\n"))
		for _, j := range column {
			name := FindString(j, []byte("name="))
			if name != nil {
				ds.Name = string(name.([]uint8))
			}
			ip := FindString(j, []byte("ip="))
			if ip != nil {
				ds.IP = string(ip.([]uint8))
			}
			password := FindString(j, []byte("password="))
			if password != nil {
				ds.Pwd = string(password.([]uint8))
			}
			weight := FindString(j, []byte("weight="))
			if weight != nil {
				s := weight.([]uint8)
				str := make([]byte, 0)
				for _, v := range s {
					if v <= 57 && v >= 48 {
						str = append(str, v)
					}
				}
				d, err := strconv.Atoi(string(str))
				if err != nil {
					log.Panicln(err)
				}
				ds.Weigh = d
			}
		}
		dn = append(dn, ds)
	}
	return dn
}
func CreateLocalPWD() []byte {
	rands := []byte("0123456789abcdefghjiklmnopqrstuvwxyz#$&*_+=")
	outs := make([]byte, 0)
	for i := 0; i < 128; i++ {
		outs = append(outs, rands[getRandomInt(0, len(rands)-1)])
	}
	f, err := os.OpenFile("./pass.ewm", os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	_, err = f.Write(outs)
	if err != nil {
		log.Println(err)
	}
	return outs
}
func Encrypt(string2 []byte) []byte {
	fs, err := os.OpenFile("./pass.ewm", os.O_RDONLY|os.O_CREATE, 666)
	if err != nil {
		log.Panicln(err)
		os.Exit(0)
	}
	defer fs.Close()
	n1, err := ioutil.ReadAll(fs)
	if len(n1) == 0 || err != nil {
		log.Panicln("There is no password in the document", err)
	}

	for k, v := range string2 {
		string2[k] = v + n1[len(n1)-1%int(v)]
	}
	return string2
}
func Decrypt(s []byte) []byte {
	fs, err := os.OpenFile("./pass.ewm", os.O_RDONLY|os.O_CREATE, 666)
	if err != nil {
		log.Println(err)
	}
	defer fs.Close()
	n1, err := ioutil.ReadAll(fs)
	for k, v := range s {
		s[k] = v - n1[len(n1)-1%int(v)]
	}
	return s
}

// save memory data to local , default 60s run one ,but you can advance or delay
func localStorageFile() {
	allkey := balala.Get("*")
	fmt.Println(allkey)
	fs, err := os.OpenFile("./spruce.db", os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		log.Println(err)
	}
	defer fs.Close()
	_, err = fs.Write(Encrypt([]byte(allkey)))
}

// 这个b方法怎么写哟 ，不球晓得，TMD
func remoteStoregeFile() {
	// 获取所有远程机器
	oAll := AllSlot
	// 饭后依次遍历 ，让其他电脑也同事备份
	for _, v := range oAll {
		go getRemote([]byte("*"), v.IP)
	}
	// 如果不出错，那么其他掉也会同时保存
}

// 这个方法是用于重置slot ，我们要重新计算 hash槽
// 实现方法 ,将每台电脑的数组取出来，重新取余，将值转移到对应slot，本机删除，如果计算结果是本机，那么不转移

func ResetSlot() {

}
func getRandomInt(start, end int) int {
	randomMutex.Lock()
	<-time.After(1 * time.Nanosecond)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := start + r.Intn(end-start+1)
	randomMutex.Unlock()
	return n
}
