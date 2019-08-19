package test

import (
	"../core/sphash"
	"testing"
)

type Data struct {
	key   rune
	value rune
}

func TestHash(t *testing.T) {
	sphash.PrepareCryptTable()
	data := sphash.HashString([]rune("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), 1)
	t.Log(data)
	data2 := sphash.HashString([]rune("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), 2)
	t.Log(data2)
	data3 := sphash.HashString([]rune("ccccccccccccccccccccccccccccccccccc"), 3)
	t.Log(data3)
	//runtime.GC()
	//var data [10000000]*Data
	//var ptr *Data
	//for i := 0; i < 10000000; i++ {
	//	ptr = &Data{int32(i), int32(i)}
	//
	//	data[i] = ptr
	//}
	//_, _ = fmt.Scan()
}
