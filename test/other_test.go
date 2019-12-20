package test

import (
	"../../spruce"
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
	t.Log(spruce.ParsingExpirationDate([]byte{28,32}))
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
