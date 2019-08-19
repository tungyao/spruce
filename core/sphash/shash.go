package sphash

type Hash struct {
	key   string
	value string
	node  *Hash
}

var cryptTable [0x500]uint

func NewHash() *Hash {
	return new(Hash)
}
func (h *Hash) Set(key string, value string) bool {
	return false
}
func (h *Hash) Get(key string) *Hash {
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
		for index2 = index1; i < 5; i++ {
			var (
				tp1 uint
				tp2 uint
			)
			seed = (seed*125 + 3) % 0x2AAAAB
			tp1 = (seed & 0xFFFF) << 0x10
			seed = (seed*125 + 3) % 0x2AAAAB
			tp2 = seed & 0xFFFF
			cryptTable[index2] = tp1 | tp2
			index2 += 0x100
		}
		i = 0

	}
}
func HashString(str []rune, hashtype uint) uint {
	var (
		key        = str
		seed1 uint = 0x7FED7FED
		seed2 uint = 0xEEEEEEEE
		ch    uint
	)
	for _, v := range key {
		va := v
		ch = uint(va)
		va++
		seed1 = cryptTable[(hashtype<<8)+ch] ^ (seed1 + seed2)
		seed2 = ch + seed1 + seed2 + (seed2 << 5) + 3

	}
	return seed1
}
func GetHashPos(str []rune) {
	var (
		HASH_OFFSET uint = 0
		HASH_A      uint = 1
		HASH_B      uint = 2
	)
	var nHash int = int(HashString(str, HASH_OFFSET))
	var nHashA int = int(HashString(str, HASH_A))
	var nHashB int = int(HashString(str, HASH_B))
	var nHashStart int = nHash % len(cryptTable)
	var nHashPos = nHashStart

}
