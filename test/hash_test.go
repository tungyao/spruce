package test

import (
	"fmt"
	"testing"
)

type Data struct {
	key   rune
	value rune
}

func TestHash(t *testing.T) {
	//sphash.PrepareCryptTable()
	//data:=sphash.HashString([]rune("username"))
	//t.Log(data)
	//runtime.GC()
	var data [10000000]*Data
	var ptr *Data
	for i := 0; i < 10000000; i++ {
		ptr = &Data{int32(i), int32(i)}

		data[i] = ptr
	}
	_, _ = fmt.Scan()
}
