package test

import (
	"../core/sphash"
	"testing"
)

func TestHash(t *testing.T) {
	sp := sphash.NewMapping()
	//sp.Set("USERA", "PASSWORD")
	//sp.Set("USERB", "PASSWORD")
	//sp.Set("USERC", "PASSWORD")
	//t.Log(sp.Get("USERC"))

	//var data [10000000]*Data
	//var ptr *Data
	d := "abcdefghij"
	for i := 0; i < len(d)*1000000; i++ {
		sp.Set(string(d[len(d)*1000000%10]), "PASSWORD")
	}
	//time.Sleep(time.Second*1)
	//_, _ = fmt.Scan()
}
