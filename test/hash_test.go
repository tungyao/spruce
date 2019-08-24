package test

import (
	"../core/sphash"
	"testing"
)

func TestHash(t *testing.T) {
	sphash.PrepareCryptTable()
	sphash.Set("USERA", "123")
	sphash.Set("USERA", "123123")
	sphash.Set("USERB", "123123123")
	//t.Log(sphash.Get("USERB"))
	//for i := 0; i < 10000000; i++ {
	//	sphash.Set(string(rand.Int()), "ab"+string(i+i))
	//}
	//for i := 0; i < 10; i++ {
	//	t.Log(sphash.Get(string(rand.Int())))
	//}
	//var data [10000000]*Data
	//var ptr *Data
	//time.Sleep(time.Second*1)
	//_, _ = fmt.Scan()
}
