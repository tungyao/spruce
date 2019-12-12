package test

import (
	"../../spruce"
	"testing"
)

func TestDIS2(t *testing.T) {
	spruce.StartSpruceDistributed(spruce.Config{
		ConfigType:    spruce.FILE,
		DCSConfigFile: "./config.yml",
		Addr:          ":88",
		KeepAlive:     false,
		IsBackup:      false,
	})
}
func TestDIS3(t *testing.T) {
	go spruce.StartSpruceDistributed(spruce.Config{})
	//a, err := net.Listen("tcp", ":79")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//for {
	//	l, _ := a.Accept()
	//	spruce.New().Set("hello", "world")
	//	t.Log(spruce.New().Get("hello"))
	//	l.Close()
	//}
}
func TestSplit(t *testing.T) {
	t.Log(spruce.SplitString([]byte("set**hello**word"), []byte("**")))
}
