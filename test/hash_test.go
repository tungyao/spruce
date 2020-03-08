package test

import (
	"../../spruce"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestHash(t *testing.T) {
	h := spruce.CreateHash(512)
	for i := 0; i < 100; i++ {
		h.Set([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(i)), 0)
	}
	//t.Log(string(h.Delete([]byte("9"))))
	////t.Log(string(h.Get([]byte("9"))))
	////h.Set([]byte(strconv.Itoa(9)), []byte(strconv.Itoa(9)),0)
	//for i := 0; i < 100; i++ {
	//	t.Log(string(h.Get([]byte(strconv.Itoa(i)))))
	//}
	//t.Log(string(h.Get([]byte(strconv.Itoa(88)))))
	h.Storage()
}
func TestReplace(t *testing.T) {
	a := "abcd\n\t\r摇动"
	fmt.Println(spruce.ReplaceTabCharacter([]byte(a)))
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
