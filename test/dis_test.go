package test

import (
	"../../spruce"
	"fmt"
	"net/http"
	"testing"
)

func TestDIS2(t *testing.T) {
	spruce.StartSpruceDistributed(spruce.Config{
		ConfigType:    spruce.FILE,
		DCSConfigFile: "./config.yml",
		Addr:          "127.0.0.1:88",
		KeepAlive:     false,
		IsBackup:      false,
		NowIP:         "127.0.0.1:88",
	})
}
func TestBack(t *testing.T) {

	http.HandleFunc("/door", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.RemoteAddr)
		writer.Write([]byte(`{"run":"yes","ip":"101.132.172.196","port":"443","continue_time":"5"}`))
	})
	http.HandleFunc("/report", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.URL.RawQuery)
	})
	http.ListenAndServe(":80", nil)
}
func TestDIS3(t *testing.T) {
	var edges = make([][]int, 0)
	edges = append(edges, []int{1, 2})
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
func TestMd5(t *testing.T) {
	t.Log(string(spruce.MD5([]byte("123"))))
}
