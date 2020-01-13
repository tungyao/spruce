package test

import (
	"../../spruce"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	h := spruce.CreateHash(512)
	h.Set([]byte("abcd"), []byte("efgh"), 2)
	fmt.Println(h.Get([]byte("abcd")))
	time.Sleep(time.Second * 3)
	fmt.Println(h.Get([]byte("abcd")))
	h.Set([]byte("abcd"), []byte("efgh"), 2)
	fmt.Println(h.Get([]byte("abcd")))

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
