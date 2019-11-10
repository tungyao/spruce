package test

import (
	"../../spruce"
	"math/rand"
	"strconv"
	"testing"
)

func TestHash(t *testing.T) {
	h := spruce.CreateHash(512)
	for i := 0; i < 100000; i++ {
		h.Set(strconv.Itoa(i), strconv.Itoa(i), 7000)
	}
	for i := 0; i < 100; i++ {
		h.Delete(strconv.Itoa(i))
	}
	for i := 0; i < 100000; i++ {
		if h.Get(strconv.Itoa(i)) == "" {
			t.Log(i)
		}
	}
}
func TestGHash(t *testing.T) {
	hash := make(map[string]string, 10240)
	for i := 0; i < 100000; i++ {
		hash[strconv.Itoa(i+rand.Int())] = strconv.Itoa(i)
	}
	//for i := 0; i < 100000; i++ {
	//	t.Log(hash[strconv.Itoa(i)+"abc"])
	//}
}
