package sphash

type Hash struct {
	key   string
	value string
	node  *Hash
}

var cryptTable [0x500]int

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
		seed   int = 0x00100001
		index1     = 0
		index2     = 0
		i          = 0
	)
	for index1 = 1; index1 < 0x100; index1++ {
		for index2 = index1; i < 5; i++ {
			var (
				tp1 int
				tp2 int
			)
			seed = (seed*125 + 3) % 0x2AAAAB
			tp1 = (seed & 0xFFFF) << 0x10
			seed = (seed*125 + 3) % 0x2AAAAB
			tp2 = seed & 0xFFFF
			cryptTable[index2] = tp1 | tp2
			index2 += 0x100
		}
	}
}
func HashString(str []rune) int {
	var (
		key   = str
		seed1 = 0x7FED7FED
		seed2 = 0xEEEEEEEE
		ch    int
	)
	for _, v := range key {
		va := v
		va++

		ch = int(va)
		seed1 = cryptTable[(1<<8)+ch] ^ (seed1 + seed2)
		seed2 = ch + seed1 + seed2 + (seed2 << 5) + 3

	}
	return seed1
}
