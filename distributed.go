package spruce

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"src/github.com/tungyao/ymload"
	"strconv"
	"sync"
	"time"
)

const (
	FILE = iota
	COMMAND
)

type Config struct {
	ConfigType    int    `配置方式`
	DCSConfigFile string `分布式的配置文件路径`
	Addr          string `跑在哪个端口上`
	NowIP         string `当前服务器运行的IP地址 暂时必须`
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
func StartSpruceDistributed(config Config) {
	CheckConfig(&config, Config{
		ConfigType:    FILE,
		DCSConfigFile: "./spruce.yml",
		Addr:          ":9102",
	})
	//fmt.Println(config)
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

// TODO Client
func client(config Config) {
	slot := initDNode(config.DCSConfigFile)
	slot.Face.IP = config.NowIP
	//if config.KeepAlive {
	//	tcpAddr, err := net.ResolveTCPAddr("tcp", config.Addr) //创建 tcpAddr数据
	//	a, err := net.ListenTCP("tcp", tcpAddr)
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//	defer a.Close()
	//	fmt.Println("\n\nserver is listening =>", a.Addr().String())
	//	fmt.Println("server is running   =>", os.Getpid())
	//	for {
	//		c, err := a.AcceptTCP()
	//		if err != nil {
	//			log.Println(err)
	//			_ = c.Close()
	//		}
	//		go func(c *net.TCPConn) {
	//			err = c.SetKeepAlive(true)
	//			err = c.SetKeepAlivePeriod(time.Second * 10)
	//			if err != nil {
	//				log.Println(err)
	//				err = c.Close()
	//				if err != nil {
	//					log.Println(err)
	//				}
	//			}
	//			data := make([]byte, 1024)
	//			n, err := c.Read(data)
	//			if err != nil {
	//				log.Println(err)
	//			} else {
	//				err = c.CloseRead()
	//			}
	//			//fmt.Println(string(data))
	//			str := SplitString(data[:n], []byte("*$"))
	//			msg := make([]byte, 0)
	//			switch string(str[0]) {
	//			case "get":
	//				msg = slot.Get(str[1])
	//			case "set":
	//				if len(str) == 3 {
	//					slot.Set(string(str[1]), string(str[2]), 0)
	//				} else if len(str) == 4 {
	//					ns, err := strconv.Atoi(string(str[3]))
	//					if err == nil {
	//						slot.Set(string(str[1]), string(str[2]), ns)
	//					} else {
	//					}
	//				}
	//			default:
	//				msg = []byte{0x65, 0x72, 0x72, 0x6F, 0x72}
	//			}
	//			_, err = c.Write([]byte(msg))
	//			if err != nil {
	//				log.Println(err)
	//			} else {
	//				err = c.CloseWrite()
	//			}
	//			//err = c.Close()
	//			//if err != nil {
	//			//	log.Println(err)
	//			//}
	//		}(c)
	//	}
	//}
	//tcpAddr, err := net.ResolveTCPAddr("tcp", config.Addr) //创建 tcpAddr数据
	a, err := net.Listen("tcp", config.Addr)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("now ip is ", config.NowIP, "we would contrast it")
	fmt.Println("\n\nserver is listening =>", a.Addr().String())
	fmt.Println("server is running   =>", os.Getpid())
	go listenAllSlotAction()
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
			case 1:
				// TODO 写到这里来了
				msg = slot.Set(data[:n])
				//return SendStatusMessage()
			case 2:
				msg = slot.Get(data[11:])
			case 3:
			}
			//switch string(str[0]) {
			//case "get":
			//	msg = slot.Get(str[1])
			//case "set":
			//	if len(str) == 3 {
			//		slot.Set(string(str[1]), string(str[2]), 0)
			//		//msg = "1"
			//	} else if len(str) == 4 {
			//		ns, err := strconv.Atoi(string(str[3]))
			//		if err == nil {
			//			slot.Set(string(str[1]), string(str[2]), ns)
			//			//msg = "1"
			//		} else {
			//			//msg = ""
			//		}
			//	}
			//case "reset":
			//	rWeigh, err := strconv.Atoi(string(str[3]))
			//	if err != nil {
			//		log.Println(err)
			//	} else {
			//		rName := string(str[1])
			//		rIp := string(str[2])
			//		rPwd := string(str[4])
			//		p := DNode{
			//			Name:  rName,
			//			IP:    rIp,
			//			Weigh: rWeigh,
			//			Pwd:   rPwd,
			//		}
			//		AllSlot = append(AllSlot, p)
			//		action <- 1
			//	}
			//default:
			//	msg = []byte{0x65, 0x72, 0x72, 0x6F, 0x72}
			//}
			_, err = c.Write(msg)
			if err != nil {
				log.Println(err)
			}
			err = c.Close()
		}(c)

	}
	//}

}
func (s *Slot) Get(key []byte) []byte {
	n := s.getHashPos(key)
	fmt.Println("get value of", n, "slot", key)
	if s.All[n].IP == s.Face.IP {
		return balala.Get(key)
	} else {
		return getRemote(SendGetMessage(key), s.All[n].IP)
	}
}
func (s *Slot) Position(key []byte) int {
	return int(s.getHashPos(key))
}
func (s *Slot) Set(lang []byte) []byte {
	key, value := SplitKeyValue(lang[11:])
	ns := s.getHashPos(key)
	fmt.Println("set value to", s.Face.IP, "slot", key)
	if s.All[ns].IP == s.Face.IP {
		fmt.Println("save")
		ti := ParsingExpirationDate(lang[2:3]).(int64)
		it := balala.Set(key, value, ti)
		return []byte{uint8(it)}
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

// save memory data to local , default 60s run one ,but you can advance or delay
func localStorageFile() {
	allkey := balala.Get([]byte("*"))
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

// 每一次触发都会使指针向后面移动
type MessageRing struct {
	Msg [2048]*msg `用来存放没有成功推送出去的消息,消息都会及时推送出去`
	All int
	mux sync.Mutex
	P   int
	Tag bool
}
type Ring struct {
	Rings []*MessageRing
	N     int
	mux   sync.Mutex
}
type msg struct {
	From    string
	To      string
	Time    string
	Content string
	Repeat  int `重复发送次数 ，超过3次就删除`
	Loop    int `在环上是第几圈`
}

// TODO Message
// 要 new 出一个环形的msg
func NewMessage() *Ring {
	nMessageRing = new(MessageRing)
	pMsg = make(chan msg, 1000)
	return &Ring{
		Rings: make([]*MessageRing, 1000),
	}
}
func (m *Ring) MSG(from, to, times, content string) msg {
	return msg{
		From:    from,
		To:      to,
		Time:    times,
		Content: content,
		Repeat:  0,
		Loop:    0,
	}
}
func (m *Ring) Push(sg msg) {
	if m.N == 0 {
		m.Rings[0] = &MessageRing{
			Msg: [2048]*msg{},
			All: 0,
			mux: sync.Mutex{},
			P:   0,
			Tag: false,
		}
		m.N++
	}
	for i := 0; i < m.N; i++ {
		m.mux.Lock()
		if m.Rings[i].Tag {
			continue
		}
		p := hashString([]rune(sg.To)) % len(m.Rings[i].Msg)
		m.Rings[i].mux.Lock()
		if m.Rings[i].Msg[p] != nil {
			mr := new(MessageRing)
			mr.Msg[p] = &sg
			mr.All++
			m.N++
			m.Rings = append(m.Rings, mr)
			break
		}
		if m.Rings[i].Msg[p] == nil {
			m.Rings[i].Msg[p] = &sg
		} else {
			continue
		}
		m.Rings[i].mux.Unlock()
		m.mux.Unlock()
		break
	}

}
func (m *Ring) Pull() {
}

// 用来表示当前的在环的那个位置
func (m *MessageRing) Pointer() int {
	return m.P
}

// 继续循环或者开始循环
func (m *MessageRing) Loop() {
}
func hashString(str []rune) int {
	var (
		key        = str
		seed1 uint = 0x7FED7FED
		seed2 uint = 0xEEEEEEEE
		ch    uint = 0
	)
	for _, v := range key {
		ch = uint(v)
		seed1 = CRY[(2<<2)+ch] ^ (seed1 + seed2)
		seed2 = ch + seed1 + seed2 + (seed2 << 5) + 3
	}
	return int(seed1)
}
