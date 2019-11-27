package spruce

import (
	"time"
)

type node struct {
	key   string
	value string
	at    int64 `insertion time -> unix time -> second`
	et    int64 `Expiration time -> second`
	next  *node
	deep  int
}
type hash struct {
	ver   []*node
	clone int
}

var CRY []uint

func CreateHash(n int) *hash {
	cryptTable := make([]uint, n)
	verticalTable := make([]*node, n)
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
		seed = (seed*uint(len(verticalTable)) + 3) % 0x2AAAAB
		tp1 = (seed & 0xFFFF) << 0x10
		seed = (seed*uint(len(verticalTable)) + 3) % 0x2AAAAB
		tp2 = seed & 0xFFFF
		cryptTable[index1] = tp1 | tp2
		i += 1
	}
	CRY = cryptTable
	return &hash{
		ver: verticalTable,
	}

}
func find(key string, node *node) string {
	tmp := node
	if tmp == nil || time.Now().Unix()-tmp.at > tmp.et && tmp.et != 0 {
		return ""
	}
	for tmp != nil {
		if tmp.key != key {
			tmp = tmp.next
			continue
		}
		break
	}
	if time.Now().Unix()-tmp.at > tmp.et && tmp.et != 0 {
		return ""
	}
	return tmp.value
}
func newNode(k, v string, deep int, exptime int64) *node {
	return &node{
		key:   k,
		value: v,
		at:    time.Now().Unix(),
		et:    exptime,
		deep:  deep,
	}
}
func (h *hash) Set(key string, value string, expTime int64) int {
	pos := h.GetHashPos([]rune(key))
	d := h.ver[pos]
	if d == nil {
		h.ver[pos] = newNode(key, value, 0, expTime)
		return int(pos)
	}
	if d.key == key {
		h.ver[pos].value = value
		return int(pos)
	}
	for d.next != nil {
		d = d.next
	}
	h.clone += 1
	d.next = newNode(key, value, d.deep+1, expTime)
	return int(pos)
}
func (h *hash) Get(key string) string {
	if key == "*" {
		return findAll(h.ver, 1)
	}
	pos := h.GetHashPos([]rune(key))
	return find(key, h.ver[pos])
}
func findAll(n []*node, tp int) string {
	tmp := n
	s := ""
	for _, v := range tmp {
		t := v
		for t != nil {
			s += t.key + "****" + t.value + "\n"
			//s += t.key + "\t" + t.value + "\n"
			t = v.next
		}
	}
	return ""
}
func (h *hash) hashString(str []rune, hashcode uint) uint {
	var (
		key        = str
		seed1 uint = 0x7FED7FED
		seed2 uint = 0xEEEEEEEE
		ch    uint = 0
	)
	for _, v := range key {
		ch = uint(v)
		seed1 = CRY[(hashcode<<2)+ch] ^ (seed1 + seed2)
		seed2 = ch + seed1 + seed2 + (seed2 << 5) + 3
	}
	return seed1
}
func (h *hash) GetHashPos(str []rune) uint {
	var (
		hOffset  uint = 0
		hashA    uint = 1
		hashB    uint = 2
		nHash    uint = h.hashString(str, hOffset)
		nHashA   uint = h.hashString(str, hashA)
		nHashB   uint = h.hashString(str, hashB)
		nHashPos uint = (nHash + nHashA + nHashB) % uint(len(CRY))
	)
	return nHashPos
}
func (h *hash) Clone() int {
	return h.clone
}
