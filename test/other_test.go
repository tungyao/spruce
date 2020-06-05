package test

import (
	"../../spruce"
	"fmt"
	"testing"
	"time"
)

func TestOther(t *testing.T) {
	s, _ := spruce.Encrypt([]byte("hello wolrd"))
	t.Log(string(spruce.Decrypt(s)))
	//t.Log(string(spruce.CreateLocalPWD()))
}
func TestStringToInt(t *testing.T) {
	//a := "12"
	//t.Log(StringToInt(a))
	//spruce.SendSetMessage([]byte("hello"),"")
}
func TestByte(t *testing.T) {
	t.Log(spruce.ParsingExpirationDate([]byte{28, 32}))
	t.Log(spruce.ParsingExpirationDate(7200))
}
func StringToInt(a string) int {
	var intSize int
	var isF bool
	slen := len(a)
	if 0 < slen && slen < 10 {
		intSize = 32
	} else if 0 < slen && slen < 19 {
		intSize = 64
	}
	if a[0] == '-' {
		isF = true
	}
	cutoff := uint64(1 << uint(intSize-1))
	println(cutoff)
	if isF {
		return -int(cutoff)
	} else {
		return int(cutoff - 1)
	}
}
func TestChannel(t *testing.T) {
	cha := make(chan int, 10)
	go func(c chan int) {
		for v := range c {
			fmt.Println(v)
		}
	}(cha)
	go func(c chan int) {
		select {
		case cha <- 0:
			fmt.Println("put key")
		}
	}(cha)
	for i := 0; i < 100; i++ {
		cha <- i
	}
}
func TestAfter(t *testing.T) {
	for {
		fmt.Println(1)
		<-time.After(time.Second * 1)
	}
}
func init() {
	fmt.Println(1)
}
func init() {
	fmt.Println(2)
}
func init() {
	fmt.Println(3)
}
func TestOne(t *testing.T) {
	d := []int{1, 2, 3, 5}
	fmt.Println(d[:2])
	fmt.Println(d[2:])
	var s []int
	s = append(s, 1)
	var m map[string]int
	m["one"] = 1
}
