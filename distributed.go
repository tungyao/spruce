package spruce

import (
	"fmt"
	"io"
	"log"
	"net"
)

var (
	slot   int
	balala = CreateHash(4096)
)

const (
	FILE = iota
	COMMAND
)

type Config struct {
	ConfigType    int    `配置方式`
	DCSConfigFile string `分布式的配置文件路径`
	Test          string
}
type DNode struct {
	Name  string
	IP    string
	Weigh int
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
	return &Slot{Count: len(c), Face: c[0], Other: c[1:], All: c, cry: cryptTable}
}
func New() *Slot {
	n := []DNode{{
		Name:  "",
		IP:    "127.0.0.1:80",
		Weigh: 1,
	}, {
		Name:  "",
		IP:    "127.0.0.1:81",
		Weigh: 2,
	}, {
		Name:  "",
		IP:    "127.0.0.1:83",
		Weigh: 3,
	}}
	return setAllDNode(n)

}
func StartSpruceDistributed(config Config) {
	CheckConfig(&config, Config{
		ConfigType:    FILE,
		DCSConfigFile: "./candi.json",
	})
	client(config.Test)
}
func client(addr string) {
	a, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		c, err := a.Accept()
		if err != nil {
			log.Println(err)
		}
		_, err = c.Write([]byte(balala.Get(string(GetData(c)))))
		if err != nil {
			log.Println(err)
		}
		err = c.Close()
	}
}
func (s *Slot) Get(key string) string {
	n := s.getHashPos([]rune(key))
	if s.All[n].IP == s.Face.IP {
		return balala.Get(key)
	} else {
		fmt.Println(s.All[n].IP)
		return getRemote(s.All[n].IP, key)
	}
}
func getRemote(ip string, key string) string {
	con, err := net.Dial("tcp", ip)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer con.Close()
	con.Write([]byte(key))
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
	log.Println(string(o))
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
