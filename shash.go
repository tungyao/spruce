package spruce

import (
	"fmt"
	"sync"
	"time"
)

type node struct {
	key   []byte
	value []byte
	at    int64 `insertion time -> unix time -> second`
	et    int64 `Expiration time -> second`
	next  *node
	deep  int
	check bool // 用来检测当前插槽是不是有值存在
	dl    int8
}
type Hash struct {
	ver   []*node
	clone int
	rw    sync.RWMutex
}

var CRY []uint

func CreateHash(n int) *Hash {
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
	return &Hash{
		ver: verticalTable,
	}

}
func find(key []byte, node *node) []byte {
	tmp := node
	if tmp == nil || (time.Now().Unix()-tmp.at > tmp.et && tmp.et != 0) {
		//tmp.check = false
		return nil
	}
	for tmp != nil {
		if !Equal(tmp.key, key) {
			tmp = tmp.next
			continue
		}
		break
	}
	fmt.Println(&tmp)
	if tmp == nil || time.Now().Unix()-tmp.at > tmp.et && tmp.et != 0 {
		//tmp = tmp.next
		tmp.check = false
		return nil
	}
	return tmp.value
}
func newNode(k, v []byte, deep int, exptime int64) *node {
	return &node{
		key:   k,
		value: v,
		at:    time.Now().Unix(),
		et:    exptime,
		deep:  deep,
		check: true,
		dl:    0,
	}
}
func (h *Hash) Set(key []byte, value []byte, expTime int64) int {
	pos := h.GetHashPos(key)
	d := h.ver[pos]
	if d == nil {
		h.ver[pos] = &node{
			key:   key,
			value: value,
			at:    time.Now().Unix(),
			et:    expTime,
			next:  nil,
			deep:  0,
			check: true,
			dl:    0,
		}
		return int(pos)
	}
	if Equal(d.key, key) {
		h.ver[pos] = &node{
			key:   key,
			value: value,
			at:    time.Now().Unix(),
			et:    expTime,
			next:  d.next,
			deep:  0,
			check: true,
			dl:    0,
		}
		return int(pos)
	}
	for d.next.check == true {
		d = d.next
	}
	h.clone += 1
	d = &node{
		key:   key,
		value: value,
		at:    time.Now().Unix(),
		et:    expTime,
		next:  d.next,
		deep:  0,
		check: true,
		dl:    0,
	}
	return int(pos)
}
func (h *Hash) Get(key []byte) []byte {
	pos := h.GetHashPos(key)
	return find(key, h.ver[pos])
}
func (h *Hash) Storage() {
	h.rw.RLock()
	defer h.rw.RUnlock()
	FindAll(h.ver)
}

// 重新从文件中读取到内存中来
func (h *Hash) Reload() {
	h.rw.RLock()
	defer h.rw.RUnlock()
}

// 这个直接读取是程序启动时会默认执行的
func (h *Hash) Load() {
	h.rw.RLock()
	defer h.rw.RUnlock()
}
func (h *Hash) Delete(key []byte) []byte {
	pos := h.GetHashPos(key)
	n, v := delete(key, h.ver[pos])
	h.ver[pos] = n
	return v
}
func (h *Hash) hashString(str []byte, hashcode uint) uint {
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
func (h *Hash) GetHashPos(str []byte) uint {
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
func (h *Hash) Clone() int {
	return h.clone
}

func delete(key []byte, nod *node) (*node, []byte) {
	if nod == nil {
		return nod, nil
	}
	p1 := nod
	p2 := nod.next
	for p2 != nil {
		if Equal(p2.key, key) {
			p1.next = p2.next
			p2 = &node{}
			//return p1, v
		} else {
			p1 = p2
		}
		p2 = p1.next
	}
	return p1, nil
}
func FindAll(n []*node) []byte {
	tmp := n
	data := make([]byte, 0)
	for _, v := range tmp {
		t := v
		for t != nil {
			if len(t.key) != 0 {
				data = append(data)
				fmt.Println(string(t.key), string(t.value))
			}
			t = t.next
		}
	}
	return nil
}

// replace tab character function
// \n -> 0
// \r -> 1
// \t -> 2
func ReplaceTabCharacter(in []byte) []byte {
	for k, v := range in {
		switch v {
		case '\n':
			in[k] = 0x0
		case '\r':
			in[k] = 0x1
		case '\t':
			in[k] = 0x2
		}
	}
	return in
}
func ReplaceTabCharacterToNormal(in []byte) []byte {
	for k, v := range in {
		switch v {
		case 0x0:
			in[k] = '\n'
		case 0x1:
			in[k] = '\r'
		case 0x2:
			in[k] = '\t'
		}
	}
	return in
}
