package test

import (
	"../../spruce"
	"fmt"
	"testing"
)

func TestOther(t *testing.T) {
	s := spruce.Encrypt([]byte("hello wolrd"))
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
func removeDuplicates(nums []int) int {
	in := make(map[int]int)
	out := make([]int, 0)
	for _, v := range nums {
		if in[v] == 0 {
			in[v] = 1
			out = append(out, v)
		}
	}
	return out
}
