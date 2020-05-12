package spruce

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type node struct {
	Key   []byte
	Value interface{}
	at    int64 `insertion time -> unix time -> second`
	et    int64 `Expiration time -> second`
	next  *node
	deep  int
	check bool // 用来检测当前插槽是不是有值存在
	dl    int8
	mux   sync.RWMutex
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
func find(key []byte, node *node) interface{} {
	tmp := node
	if tmp == nil || (time.Now().Unix()-tmp.at > tmp.et && tmp.et != 0) {
		// tmp.check = false
		return nil
	}
	for tmp != nil {
		if !Equal(tmp.Key, key) {
			tmp = tmp.next
			continue
		}
		break
	}
	if tmp == nil || time.Now().Unix()-tmp.at > tmp.et && tmp.et != 0 {
		tmp.check = false
		return nil
	}
	if tmp.check == false {
		return nil
	}
	return tmp.Value
}
func newNode(k, v []byte, deep int, exptime int64) *node {
	return &node{
		Key:   k,
		Value: v,
		at:    time.Now().Unix(),
		et:    exptime,
		deep:  deep,
		check: true,
		dl:    0,
	}
}
func (h *Hash) Set(key []byte, value interface{}, expTime int64) int {
	h.rw.RLock()
	defer h.rw.RUnlock()
	pos := h.GetHashPos(key)
	d := h.ver[pos]
	if d == nil {
		h.ver[pos] = &node{
			Key:   key,
			Value: value,
			at:    time.Now().Unix(),
			et:    expTime,
			next:  nil,
			deep:  0,
			check: true,
			dl:    0,
		}
		return int(pos)
	}
	if Equal(d.Key, key) {
		h.ver[pos] = &node{
			Key:   key,
			Value: value,
			at:    time.Now().Unix(),
			et:    expTime,
			next:  d.next,
			deep:  0,
			check: true,
			dl:    0,
		}
		return int(pos)
	}
	for d.next != nil && d.next.check == true {
		d = d.next
	}
	h.clone += 1
	d = &node{
		Key:   key,
		Value: value,
		at:    time.Now().Unix(),
		et:    expTime,
		next:  d.next,
		deep:  0,
		check: true,
		dl:    0,
	}
	return int(pos)
}
func (h *Hash) Get(key []byte) interface{} {
	pos := h.GetHashPos(key)
	if Equal(key, []byte("all")) {
		return FindAll(h.ver)
	}
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
func (h *Hash) Delete(key []byte) []byte {
	pos := h.GetHashPos(key)
	h.rw.RLock()
	_, v := delete(key, h.ver[pos])
	h.rw.RUnlock()
	// h.ver[pos] = n
	return v
}
func delete(key []byte, nod *node) (*node, []byte) {
	for nod != nil {
		p1 := nod.next
		if Equal(nod.Key, key) {
			nod.check = false
			return nod, nil
		}
		nod = p1.next
	}
	return nod, nil
}
func FindAll(n []*node) []byte {
	tmp := n
	data := ""
	for _, v := range tmp {
		t := v
		for t != nil {
			if len(t.Key) != 0 {
				// x := make([]byte, 0)
				data += string(t.Key) + "\t\t" + fmt.Sprintf("%s", t.Value) + "\n"
			}
			t = t.next
		}
	}
	return []byte(data)
}
func (h *Hash) GetAll() []interface{} {
	tmp := h.ver
	data := make([]interface{}, 0)
	for _, v := range tmp {
		t := v
		for t != nil {
			if len(t.Key) != 0 && t.check {
				data = append(data, t.Value)
			}
			t = t.next
		}
	}
	return data
}

func (h *Hash) GetAllWithKey() []*node {
	tmp := h.ver
	data := make([]*node, 0)
	for _, v := range tmp {
		t := v
		for t != nil {
			if len(t.Key) != 0 && t.check {
				data = append(data, t)
			}
			t = t.next
		}
	}
	return data
}

// set all thing to bytes
func ToBytes(x interface{}) ([]byte, error) {
	fmt.Println(reflect.TypeOf(x).Kind())
	switch reflect.TypeOf(x).Kind() {
	case reflect.String:
		return []byte(x.(string)), nil
	case reflect.Slice:
		return x.([]byte), nil
	case reflect.Struct:
		return json.Marshal(x)
	case reflect.Ptr:

	}
	return nil, errors.New("have no matched")
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
