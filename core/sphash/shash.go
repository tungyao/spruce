package sphash

import (
	"strings"
)

type Node struct {
	key   string
	value string
	next  *Node
}
type Hash struct {
	isExists bool
	hsa      uint
	hsb      uint
	deep     int
	start    *Node
	end      *Node
}

var cryptTable [0x500]uint   //TODO 用于算出hash值得数组
var verticalTable [256]*Hash //TODO 用于存放hash 取余过后 得数值
func IsEmpty(node *Node) bool {
	return node == nil
}
func IsLast(node *Node) bool {
	return node.next == nil
}
func FindPrevious(key string, value string, node *Node) *Node {
	tmp := node
	for tmp.next != nil && tmp.next.key != key {
		tmp = tmp.next
	}
	return tmp
}
func Find(key string, node *Node) *Node {
	tmp := node
	for tmp.key != key {
		tmp = tmp.next
	}
	return tmp
}
func NewNode(k, v string, deep int) *Node {
	return &Node{
		key:   k,
		value: v,
	}
}
func Set(key string, value string) {
	obj, _ := GetHashPos([]rune(key))

	newNode := NewNode(key, value, obj.deep+1)
	if obj.deep == 0 {
		obj.start = newNode
		obj.end = newNode
	} else {
		lastPost := obj.end
		lastPost.next = newNode
		obj.end = newNode
	}
	obj.deep++
	//TODO 去这个网站 https://studygolang.com/articles/12686
}
func Get(key string) *Node {
	return nil
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
func GetHashPos(str []rune) (Hash, uint) { //TODO 需要将计算过的hash值取余放入数组中
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
				return *verticalTable[nHashPos], nHashPos
			} else {
				nHashPos = (nHashPos + 1) % 256

			}
			if nHashPos == nHashStart {
				return *verticalTable[nHashPos], nHashPos
			}
		}
	}
	verticalTable[nHashPos] = &Hash{
		isExists: true,
		hsa:      nHashA,
		hsb:      nHashB,
		deep:     0,
		start:    nil,
		end:      nil,
	}
	return *verticalTable[nHashPos], nHashPos

}
