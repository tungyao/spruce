package sphash

import (
	"strings"
)

type Mapping struct {
	hash  *Hash
	deep  uint
	key   string
	value string
	node  map[string]*Mapping
}
type Hash struct {
	isExists bool
	hsa      uint
	hsb      uint
	node     map[string]string
}

var cryptTable [0x500]uint   //TODO 用于算出hash值得数组
var verticalTable [256]*Hash //TODO 用于存放hash 取余过后 得数值
func NewMapping() *Mapping {
	PrepareCryptTable()
	return new(Mapping)
}
func NewNode(h *Hash, k string, v string, d uint) *Mapping {
	return &Mapping{
		hash:  h,
		deep:  d,
		key:   k,
		value: v,
		node:  make(map[string]*Mapping),
	}
}
func (m *Mapping) Set(key string, value string) {
	_, k := GetHashPos([]rune(key))
	if verticalTable[k].node == nil {
		verticalTable[k].node = make(map[string]string)
	}
	verticalTable[k].node[key] = value
}
func (m *Mapping) Get(key string) string {
	_, k := GetHashPos([]rune(key))
	return verticalTable[k].node[key]
}

func PrepareCryptTable() {
	var (
		seed   uint = 0x00100001
		index1      = 0
		index2      = 0
		i           = 0
	)
	for index1 = 0; index1 < 0x100; index1++ {
		for index2, i = index1, 0; i < 5; index2 += 0x100 {
			var (
				tp1 uint
				tp2 uint
			)
			seed = (seed*125 + 3) % 0x2AAAAB
			tp1 = (seed & 0xFFFF) << 0x10
			seed = (seed*125 + 3) % 0x2AAAAB
			tp2 = seed & 0xFFFF
			cryptTable[index2] = tp1 | tp2
			i += 1
		}
	}
}
func HashString(str []rune, hashtype uint) uint {
	var (
		key        = str
		seed1 uint = 0x7FED7FED
		seed2 uint = 0xEEEEEEEE
		ch    uint = 0
	)
	for _, v := range key {
		ch = uint(strings.ToUpper(string(v))[0])
		seed1 = cryptTable[(hashtype<<8)+ch] ^ (seed1 + seed2)
		seed2 = ch + seed1 + seed2 + (seed2 << 5) + 3
	}
	return seed1
}
func GetHashPos(str []rune) (*Hash, uint) { //TODO 需要将计算过的hash值取余放入数组中
	var (
		HashOffset uint = 0
		HashA      uint = 1
		HashB      uint = 2
		nHash      uint = HashString(str, HashOffset)
		nHashA     uint = HashString(str, HashA)
		nHashB     uint = HashString(str, HashB)
		nHashStart uint = nHash % uint(len(cryptTable))
		nHashPos        = nHashStart % 256 //TODO 相当于经过了三次Hash,最终得出位置
	)
	for i := 0; i < len(verticalTable); i++ {
		if verticalTable[nHashPos] != nil {
			if verticalTable[nHashPos].hsa == nHashA && verticalTable[nHashPos].hsb == nHashB {
				return verticalTable[nHashPos], nHashPos
			} else {
				nHashPos = (nHashPos + 1) % 256

			}
			if nHashPos == nHashStart {
				return nil, nHashPos
			}
		}
	}
	verticalTable[nHashPos] = &Hash{
		isExists: true,
		hsa:      nHashA,
		hsb:      nHashB,
	}
	return verticalTable[nHashPos], nHashPos

}
