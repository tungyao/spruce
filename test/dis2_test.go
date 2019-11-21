package test

import (
	"../../spruce"
	"fmt"
	"testing"
)

func TestDIS1(t *testing.T) {
	//spruce.StartSpruceDistributed(spruce.Config{Test: ":80"})
	a := spruce.CreateHash(512)
	a.Set("hello", "world", 0)
	a.Set("helloa", "worlda", 0)
	a.Set("hellob", "worldb", 0)
	fmt.Print(a.Get("*"))
}
func TestDASDSA(t *testing.T) {
	//pub := "HELLOA"
	//pub := "HELLOB"
	//pub := "HELLOC"
	//pri :="WORLD"
	str := "mkasdkajsdkkasdjkajsdkakjsdasdasdaaasdnjkkeksna:_!@#$%^&*()_+=-你好吗"
	for _, v := range str {
		fmt.Println(v)
	}

}
