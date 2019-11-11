package test

import (
	"../../spruce"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	h := spruce.CreateHash(spruce.Config{})
	//for i := 0; i < 100; i++ {
	//	h.Set(strconv.Itoa(i+rand.Int()), strconv.Itoa(i),7000)
	//}
	h.Set("hello", "world", 2)
	time.Sleep(time.Second * 3)
	t.Log(h.Get("hello"))
}
func TestGHash(t *testing.T) {
	hash := make(map[string]string, 10240)
	for i := 0; i < 10240; i++ {
		hash[strconv.Itoa(i+rand.Int())] = strconv.Itoa(i)
	}
	//for i := 0; i < 100000; i++ {
	//	t.Log(hash[strconv.Itoa(i)+"abc"])
	//}
}
